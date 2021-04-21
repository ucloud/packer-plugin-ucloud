package uhost

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/hashicorp/packer-plugin-sdk/acctest/testutils"
	ucloudcommon "github.com/hashicorp/packer-plugin-ucloud/builder/ucloud/common"
	"github.com/stretchr/testify/assert"
)

func TestAccBuilder_validateRegion(t *testing.T) {
	t.Parallel()

	if os.Getenv(acctest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", acctest.TestEnvVar))
		return
	}

	err := testAccPreCheck()
	if err != nil {
		t.Fatalf(err.Error())
	}

	access := &ucloudcommon.AccessConfig{Region: "cn-bj2"}
	err = access.Config()
	if err != nil {
		t.Fatalf("Error on initing UCloud AccessConfig, %s", err)
	}

	err = access.ValidateRegion("cn-sh2")
	if err != nil {
		t.Fatalf("Expected pass with valid region but failed: %s", err)
	}

	err = access.ValidateRegion("invalidRegion")
	if err == nil {
		t.Fatal("Expected failure due to invalid region but passed")
	}
}

func TestAccBuilder_basic(t *testing.T) {
	t.Parallel()
	testCase := &acctest.PluginTestCase{
		Name:     "uhost_basic_test",
		Setup:    testAccPreCheck,
		Template: testBuilderAccBasic,
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

const testBuilderAccBasic = `
{	"builders": [{
		"type": "ucloud-uhost",
		"region": "cn-bj2",
		"availability_zone": "cn-bj2-02",
		"instance_type": "n-basic-2",
		"source_image_id":"uimage-f1chxn",
		"ssh_username":"root",
		"image_name": "packer-test-basic_{{timestamp}}"
	}]
}`

func TestAccBuilder_ubuntu(t *testing.T) {
	t.Parallel()
	testCase := &acctest.PluginTestCase{
		Name:     "uhost_ubuntu_test",
		Setup:    testAccPreCheck,
		Template: testBuilderAccUbuntu,
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

const testBuilderAccUbuntu = `
{	"builders": [{
		"type": "ucloud-uhost",
		"region": "cn-bj2",
		"availability_zone": "cn-bj2-02",
		"instance_type": "n-basic-2",
		"source_image_id":"uimage-irofn4",
		"ssh_username":"ubuntu",
		"image_name": "packer-test-ubuntu_{{timestamp}}"
	}]
}`

func TestAccBuilder_regionCopy(t *testing.T) {
	t.Parallel()
	projectId := os.Getenv("UCLOUD_PROJECT_ID")
	testCase := &acctest.PluginTestCase{
		Name:     "uhost_ubuntu_test",
		Setup:    testAccPreCheck,
		Template: testBuilderAccRegionCopy(projectId),
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return checkRegionCopy(
				projectId,
				[]ucloudcommon.ImageDestination{
					{ProjectId: projectId, Region: "cn-sh2", Name: "packer-test-regionCopy-sh", Description: "test"},
				})
		},
	}
	acctest.TestPlugin(t, testCase)
}

func testBuilderAccRegionCopy(projectId string) string {
	return fmt.Sprintf(`
{
	"builders": [{
		"type": "ucloud-uhost",
		"region": "cn-bj2",
		"availability_zone": "cn-bj2-02",
		"instance_type": "n-basic-2",
		"source_image_id":"uimage-f1chxn",
		"ssh_username":"root",
		"image_name": "packer-test-regionCopy-bj",
		"image_copy_to_mappings": [{
			"project_id":  	%q,
			"region":		"cn-sh2",
			"name":			"packer-test-regionCopy-sh",
			"description": 	"test"
		}]
	}],
	"post-processors": [
		{
		  "type": "manifest"
		}
	]
}`, projectId)
}

func checkRegionCopy(projectId string, imageDst []ucloudcommon.ImageDestination) error {
	manifest, err := testutils.GetArtifact("packer-manifest.json")
	if err != nil {
		return err
	}

	destSet := ucloudcommon.NewImageInfoSet(nil)
	for _, dest := range imageDst {
		destSet.Set(ucloudcommon.ImageInfo{
			Region:    dest.Region,
			ProjectId: dest.ProjectId,
		})
	}

	id := manifest.Builds[0].ArtifactId
	ucloudImages := strings.Split(id, ",")
	for _, r := range ucloudImages {
		info := strings.Split(r, ":")
		projId := info[0]
		region := info[1]
		imageId := info[2]
		if projId == projectId && region == "cn-bj2" {
			destSet.Remove(imageId)
			continue
		}

		if destSet.Get(projId, region) == nil {
			return fmt.Errorf("project%s : region%s is not the target but found in artifacts", projId, region)
		}

		destSet.Remove(imageId)
	}

	if len(destSet.GetAll()) > 0 {
		return fmt.Errorf("the following copying targets not found in corresponding artifacts : %#v", destSet.GetAll())
	}

	client, _ := testUCloudClient()
	for _, r := range ucloudImages {
		info := strings.Split(r, ":")
		projId := info[0]
		region := info[1]
		imageId := info[2]

		if projId == projectId && region == "cn-bj2" {
			continue
		}
		imageSet, err := client.DescribeImageByInfo(projId, region, imageId)
		if err != nil {
			if ucloudcommon.IsNotFoundError(err) {
				return fmt.Errorf("image %s in artifacts can not be found", imageId)
			}
			return err
		}

		if region == "cn-sh2" && imageSet.ImageName != "packer-test-regionCopy-sh" {
			return fmt.Errorf("the name of image %q in artifacts should be %s, got %s", imageId, "packer-test-regionCopy-sh", imageSet.ImageName)
		}
	}

	return nil
}

func testAccPreCheck() error {
	if v := os.Getenv("UCLOUD_PUBLIC_KEY"); v == "" {
		return fmt.Errorf("UCLOUD_PUBLIC_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("UCLOUD_PRIVATE_KEY"); v == "" {
		return fmt.Errorf("UCLOUD_PRIVATE_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("UCLOUD_PROJECT_ID"); v == "" {
		return fmt.Errorf("UCLOUD_PROJECT_ID must be set for acceptance tests")
	}
	return nil
}

func TestUCloudClientBaseUrlConfigurable(t *testing.T) {
	const url = "baseUrl"
	access := &ucloudcommon.AccessConfig{BaseUrl: url, PublicKey: "test", PrivateKey: "test"}
	client, err := access.Client()
	assert.Nil(t, err)
	assert.Equal(t, url, client.UAccountConn.Client.GetConfig().BaseUrl, "account conn's base url not configurable")
	assert.Equal(t, url, client.UHostConn.Client.GetConfig().BaseUrl, "host conn's base url not configurable")
	assert.Equal(t, url, client.UNetConn.Client.GetConfig().BaseUrl, "net conn's base url not configurable")
	assert.Equal(t, url, client.VPCConn.Client.GetConfig().BaseUrl, "vpc conn's base url not configurable")
}

func testUCloudClient() (*ucloudcommon.UCloudClient, error) {
	access := &ucloudcommon.AccessConfig{Region: "cn-bj2"}
	err := access.Config()
	if err != nil {
		return nil, err
	}
	client, err := access.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}
