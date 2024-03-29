package uhost

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"
	ucloudcommon "github.com/ucloud/packer-plugin-ucloud/builder/ucloud/common"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

type stepCreateInstance struct {
	Region        string
	Zone          string
	InstanceType  string
	InstanceName  string
	BootDiskType  string
	BootDiskSize  int
	SourceImageId string
	UsePrivateIp  bool

	EipBandwidth   int
	EipChargeMode  string
	UserData       string
	UserDataFile   string
	MinCpuPlatform string

	instanceId string
}

func (s *stepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*ucloudcommon.UCloudClient)
	conn := client.UHostConn
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Creating Instance...")
	req, err := s.buildCreateInstanceRequest(state)
	if err != nil {
		return ucloudcommon.Halt(state, err, "Error on build instance request")
	}

	resp, err := conn.CreateUHostInstance(req)
	if err != nil {
		return ucloudcommon.Halt(state, err, "Error on creating instance")
	}
	instanceId := resp.UHostIds[0]

	err = retry.Config{
		Tries: 100,
		ShouldRetry: func(err error) bool {
			return ucloudcommon.IsExpectedStateError(err)
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 2 * time.Second, MaxBackoff: 6 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		inst, err := client.DescribeUHostById(instanceId)
		if err != nil {
			return err
		}

		if inst.State == "ResizeFail" {
			return fmt.Errorf("resizing instance failed")
		}

		if inst.State == "Install Fail" {
			return fmt.Errorf("install failed")
		}

		if inst == nil || inst.State != ucloudcommon.InstanceStateRunning {
			return ucloudcommon.NewExpectedStateError("instance", instanceId)
		}

		return nil
	})

	if err != nil {
		return ucloudcommon.Halt(state, err, fmt.Sprintf("Error on waiting for instance %q to become available", instanceId))
	}

	ui.Message(fmt.Sprintf("Creating instance %q complete", instanceId))
	instance, err := client.DescribeUHostById(instanceId)
	if err != nil {
		return ucloudcommon.Halt(state, err, fmt.Sprintf("Error on reading instance when creating %q", instanceId))
	}

	s.instanceId = instanceId
	state.Put("instance", instance)
	// instance_id is the generic term used so that users can have access to the
	// instance id inside of the provisioners, used in step_provision.
	state.Put("instance_id", instanceId)

	if instance.BootDiskState != ucloudcommon.BootDiskStateNormal {
		ui.Say("Waiting for boot disk of instance initialized")
		if s.BootDiskType == "local_normal" || s.BootDiskType == "local_ssd" {
			ui.Message(fmt.Sprintf("Warning: It takes around 10 mins for boot disk initialization when `boot_disk_type` is %q", s.BootDiskType))
		}

		err = retry.Config{
			Tries: 200,
			ShouldRetry: func(err error) bool {
				return ucloudcommon.IsExpectedStateError(err)
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 2 * time.Second, MaxBackoff: 12 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			inst, err := client.DescribeUHostById(instanceId)
			if err != nil {
				return err
			}
			if inst.BootDiskState != ucloudcommon.BootDiskStateNormal {
				return ucloudcommon.NewExpectedStateError("boot_disk of instance", instanceId)
			}

			return nil
		})

		if err != nil {
			return ucloudcommon.Halt(state, err, fmt.Sprintf("Error on waiting for boot disk of instance %q initialized", instanceId))
		}

		ui.Message(fmt.Sprintf("Waiting for boot disk of instance %q initialized complete", instanceId))
	}

	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	if s.instanceId == "" {
		return
	}
	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)

	ui := state.Get("ui").(packersdk.Ui)
	ctx := context.TODO()

	if cancelled || halted {
		ui.Say("Deleting instance because of cancellation or error...")
	} else {
		ui.Say("Deleting instance...")
	}

	client := state.Get("client").(*ucloudcommon.UCloudClient)
	conn := client.UHostConn

	instance, err := client.DescribeUHostById(s.instanceId)
	if err != nil {
		if ucloudcommon.IsNotFoundError(err) {
			return
		}
		ui.Error(fmt.Sprintf("Error on reading instance when deleting %q, %s",
			s.instanceId, err.Error()))
		return
	}

	if instance.State != ucloudcommon.InstanceStateStopped {
		stopReq := conn.NewStopUHostInstanceRequest()
		stopReq.UHostId = ucloud.String(s.instanceId)
		if _, err = conn.StopUHostInstance(stopReq); err != nil {
			ui.Error(fmt.Sprintf("Error on stopping instance when deleting %q, %s",
				s.instanceId, err.Error()))
			return
		}

		err = retry.Config{
			Tries: 100,
			ShouldRetry: func(err error) bool {
				return ucloudcommon.IsExpectedStateError(err)
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 2 * time.Second, MaxBackoff: 6 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			instance, err := client.DescribeUHostById(s.instanceId)
			if err != nil {
				return err
			}

			if instance.State != ucloudcommon.InstanceStateStopped {
				return ucloudcommon.NewExpectedStateError("instance", s.instanceId)
			}

			return nil
		})

		if err != nil {
			ui.Error(fmt.Sprintf("Error on waiting for stopping instance when deleting %q, %s",
				s.instanceId, err.Error()))
			return
		}
	}

	deleteReq := conn.NewTerminateUHostInstanceRequest()
	deleteReq.UHostId = ucloud.String(s.instanceId)
	deleteReq.ReleaseUDisk = ucloud.Bool(true)
	deleteReq.ReleaseEIP = ucloud.Bool(true)

	if _, err = conn.TerminateUHostInstance(deleteReq); err != nil {
		ui.Error(fmt.Sprintf("Error on deleting instance %q, %s",
			s.instanceId, err.Error()))
		return
	}

	err = retry.Config{
		Tries:       50,
		ShouldRetry: func(err error) bool { return !ucloudcommon.IsNotFoundError(err) },
		RetryDelay:  (&retry.Backoff{InitialBackoff: 2 * time.Second, MaxBackoff: 6 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		_, err := client.DescribeUHostById(s.instanceId)
		return err
	})

	if err != nil {
		ui.Error(fmt.Sprintf("Error on waiting for instance %q to be deleted: %s",
			s.instanceId, err.Error()))
		return
	}

	ui.Message(fmt.Sprintf("Deleting instance %q complete", s.instanceId))
}

func (s *stepCreateInstance) buildCreateInstanceRequest(state multistep.StateBag) (*uhost.CreateUHostInstanceRequest, error) {
	client := state.Get("client").(*ucloudcommon.UCloudClient)
	conn := client.UHostConn
	srcImage := state.Get("source_image").(*uhost.UHostImageSet)
	config := state.Get("config").(*Config)
	connectConfig := &config.RunConfig.Comm

	var password string
	if srcImage.OsType == "Linux" {
		password = config.Comm.SSHPassword
	}

	if password == "" {
		password = fmt.Sprintf("%s%s%s",
			s.randStringFromCharSet(5, ucloudcommon.DefaultPasswordStr),
			s.randStringFromCharSet(1, ucloudcommon.DefaultPasswordSpe),
			s.randStringFromCharSet(5, ucloudcommon.DefaultPasswordNum))
		if srcImage.OsType == "Linux" {
			connectConfig.SSHPassword = password
		}
	}

	req := conn.NewCreateUHostInstanceRequest()
	t, _ := ucloudcommon.ParseInstanceType(s.InstanceType)

	req.CPU = ucloud.Int(t.CPU)
	req.Memory = ucloud.Int(t.Memory)
	req.Name = ucloud.String(s.InstanceName)
	req.LoginMode = ucloud.String("Password")
	req.Zone = ucloud.String(s.Zone)
	req.ImageId = ucloud.String(s.SourceImageId)
	req.ChargeType = ucloud.String("Dynamic")
	req.Password = ucloud.String(password)
	req.MinimalCpuPlatform = ucloud.String(s.MinCpuPlatform)
	req.MachineType = ucloud.String(strings.ToUpper(t.HostType))

	if v, ok := state.GetOk("security_group_id"); ok {
		req.SecurityGroupId = ucloud.String(v.(string))
	}

	if v, ok := state.GetOk("vpc_id"); ok {
		req.VPCId = ucloud.String(v.(string))
	}

	if v, ok := state.GetOk("subnet_id"); ok {
		req.SubnetId = ucloud.String(v.(string))
	}

	userData, err := s.getUserData(state)
	if err != nil {
		return nil, err
	}
	if userData != "" {
		req.UserData = ucloud.String(userData)
	}

	bootDisk := uhost.UHostDisk{}
	bootDisk.IsBoot = ucloud.String("true")
	bootDisk.Size = ucloud.Int(srcImage.ImageSize)
	if s.BootDiskSize > 0 {
		if s.BootDiskSize < srcImage.ImageSize {
			return nil, errors.New("boot disk size should not be smaller than image size")
		}
		bootDisk.Size = ucloud.Int(s.BootDiskSize)
	}
	bootDisk.Type = ucloud.String(ucloudcommon.BootDiskTypeMap.Convert(s.BootDiskType))

	req.Disks = append(req.Disks, bootDisk)

	if v, ok := state.GetOk("user_data"); ok {
		req.UserData = ucloud.String(base64.StdEncoding.EncodeToString([]byte(v.(string))))
	}

	if !s.UsePrivateIp {
		operatorName := ucloud.String("International")
		if strings.HasPrefix(s.Region, "cn-") {
			operatorName = ucloud.String("Bgp")
		}

		networkInterface := uhost.CreateUHostInstanceParamNetworkInterface{
			EIP: &uhost.CreateUHostInstanceParamNetworkInterfaceEIP{
				Bandwidth:    ucloud.Int(s.EipBandwidth),
				PayMode:      ucloud.String(ucloudcommon.ChargeModeMap.Convert(s.EipChargeMode)),
				OperatorName: operatorName,
			},
		}

		req.NetworkInterface = append(req.NetworkInterface, networkInterface)
	}
	return req, nil
}

func (s *stepCreateInstance) randStringFromCharSet(strlen int, charSet string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}

func (s *stepCreateInstance) getUserData(state multistep.StateBag) (string, error) {
	userData := s.UserData

	if s.UserDataFile != "" {
		data, err := ioutil.ReadFile(s.UserDataFile)
		if err != nil {
			return "", fmt.Errorf("error on reading user_data_file, %s", err)
		}

		userData = string(data)
	}

	if userData != "" {
		userData = base64.StdEncoding.EncodeToString([]byte(userData))
	}

	return userData, nil

}
