// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.
package common

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/MSOpenTech/packer-hyperv/packer/powershell/hyperv"
)


type StepUnmountFloppyDrive struct {
}

func (s *StepUnmountFloppyDrive) Run(state multistep.StateBag) multistep.StepAction {
	//driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	errorMsg := "Error Unmounting floppy drive: %s"
	vmName := state.Get("vmName").(string)

	ui.Say("Unmounting floppy drive (Run)...")

	err := hyperv.UnmountFloppyDrive(vmName)
	if err != nil {
		err := fmt.Errorf(errorMsg, err)
		state.Put("error", err)
		ui.Error(err.Error())
	}

	return multistep.ActionContinue
}

func (s *StepUnmountFloppyDrive) Cleanup(state multistep.StateBag) {
	// do nothing
}
