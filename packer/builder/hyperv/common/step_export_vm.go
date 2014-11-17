// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.
package common

import (
	"fmt"
	"path/filepath"
	"io/ioutil"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	powershell "github.com/MSOpenTech/packer-hyperv/packer/powershell"
)

const(
	vhdDir string = "Virtual Hard Disks"
	vmDir string = "Virtual Machines"
)

type StepExportVm struct {
	OutputDir string
}

func (s *StepExportVm) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	var err error
	var errorMsg string

	vmName := state.Get("vmName").(string)
	tmpPath :=	state.Get("packerTempDir").(string)
	outputPath := s.OutputDir

	// create temp path to export vm
	errorMsg = "Error creating temp export path: %s"
	vmExportPath , err := ioutil.TempDir(tmpPath, "export")
	if err != nil {
		err := fmt.Errorf(errorMsg, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Exporting vm...")

	var script ScriptBuilder
	script.WriteLine("param([string]$vmName, [string]$path)")
	script.WriteLine("Export-VM -Name $vmName -Path $path")

	powershell := new(powershell.PowerShellCmd)
	err = powershell.RunFile(script.Bytes(), vmName, vmExportPath)

	if err != nil {
		errorMsg = "Error exporting vm: %s"
		err := fmt.Errorf(errorMsg, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// copy to output dir
	expPath := filepath.Join(vmExportPath,vmName)

	ui.Say("Coping to output dir...")

	script.Reset()
	script.WriteLine("param([string]$srcPath, [string]$dstPath, [string]$vhdDirName, [string]$vmDir)")
	script.WriteLine("Copy-Item -Path $srcPath/$vhdDirName -Destination $dstPath -recurse")
	script.WriteLine("Copy-Item -Path $srcPath/$vmDir -Destination $dstPath")
	script.WriteLine("Copy-Item -Path $srcPath/$vmDir/*.xml -Destination $dstPath/$vmDir")

	err = powershell.RunFile(script.Bytes(), expPath, outputPath, vhdDir, vmDir)
	if err != nil {
		errorMsg = "Error exporting vm: %s"
		err := fmt.Errorf(errorMsg, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepExportVm) Cleanup(state multistep.StateBag) {
	// do nothing
}
