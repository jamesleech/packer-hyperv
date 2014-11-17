// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.
package common

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	powershell "github.com/MSOpenTech/packer-hyperv/packer/powershell"
)

type StepEnableIntegrationService struct {
	name string
}

func (s *StepEnableIntegrationService) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	vmName := state.Get("vmName").(string)
	s.name = "Guest Service Interface"

	ui.Say("Enabling Integration Service...")

	powershell := new(powershell.PowerShellCmd)

	var script ScriptBuilder
	script.WriteLine("param([string]$vmName, [string]$integrationServiceName)")
	script.WriteLine("Enable-VMIntegrationService -VMName $vmName -Name $integrationServiceName")

	err := powershell.RunFile(script.Bytes(), vmName, s.name)

	if err != nil {
		err := fmt.Errorf("Error enabling Integration Service: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepEnableIntegrationService) Cleanup(state multistep.StateBag) {
	// do nothing
}
