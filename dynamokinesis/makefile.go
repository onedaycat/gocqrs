// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/plimble/mage/sh"
)

type Build mg.Namespace

func (Build) Linux() {
	sh.BuildLinux(".", "./bin/app")
	fmt.Println("Build Done")
}

func (Build) Mac() {
	sh.BuildMac(".", "./bin/app")
	fmt.Println("Build Done")
}

func Deploy() {
	Build{}.Linux()
	sh.Exec("serverless deploy -v")
}

func Remove() {
	sh.Exec("serverless remove -v")
}
