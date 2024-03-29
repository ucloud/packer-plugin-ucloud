package uhost

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	ucloudcommon "github.com/ucloud/packer-plugin-ucloud/builder/ucloud/common"
)

type stepPreValidate struct {
	ProjectId         string
	Region            string
	Zone              string
	ImageDestinations []ucloudcommon.ImageDestination
}

func (s *stepPreValidate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if err := s.validateProjectIds(state); err != nil {
		return ucloudcommon.Halt(state, err, "")
	}

	if err := s.validateRegions(state); err != nil {
		return ucloudcommon.Halt(state, err, "")
	}

	if err := s.validateZones(state); err != nil {
		return ucloudcommon.Halt(state, err, "")
	}

	return multistep.ActionContinue
}

func (s *stepPreValidate) validateProjectIds(state multistep.StateBag) error {
	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("config").(*Config)

	ui.Say("Validating project_id and copied project_ids...")

	var errs *packersdk.MultiError
	if err := config.ValidateProjectId(s.ProjectId); err != nil {
		errs = packersdk.MultiErrorAppend(errs, err)
	}

	for _, imageDestination := range s.ImageDestinations {
		if err := config.ValidateProjectId(imageDestination.ProjectId); err != nil {
			errs = packersdk.MultiErrorAppend(errs, err)
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (s *stepPreValidate) validateRegions(state multistep.StateBag) error {
	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("config").(*Config)

	ui.Say("Validating region and copied regions...")

	var errs *packersdk.MultiError
	if err := config.ValidateRegion(s.Region); err != nil {
		errs = packersdk.MultiErrorAppend(errs, err)
	}
	for _, imageDestination := range s.ImageDestinations {
		if err := config.ValidateRegion(imageDestination.Region); err != nil {
			errs = packersdk.MultiErrorAppend(errs, err)
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (s *stepPreValidate) validateZones(state multistep.StateBag) error {
	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("config").(*Config)

	ui.Say("Validating availability_zone...")

	var errs *packersdk.MultiError
	if err := config.ValidateZone(s.Region, s.Zone); err != nil {
		errs = packersdk.MultiErrorAppend(errs, err)
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (s *stepPreValidate) Cleanup(multistep.StateBag) {}
