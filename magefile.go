//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"log"
	"os"
	"os/exec"
	"sheazuzu/build"
	"sheazuzu/sheazuzu"
)

// Do not run this inside of IntelliJ
// It will not work, because IntelliJ sets environment vars, that make the go cli expect a go.mod

var (
	VERSION = ""
	MODULE  = "main"
)

func init() {

	props, err := build.ReadProperties("versions.properties")
	if err != nil {
		log.Fatalf("Failed to read the versions.properties file: %s", err.Error())
	}
	VERSION = props["sheazuzu.version"]
}

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

func Version() {
	fmt.Println(VERSION)
}

func Clean() {
	mg.Deps(sheazuzu.Clean)
}

func Prepare() {
	mg.Deps(sheazuzu.Prepare)
}

func Generate() {
	Prepare()
	mg.Deps(sheazuzu.GenerateServer)
}

func Build() {
	Generate()
	mg.Deps(sheazuzu.Build)
}

func Run() error {
	cmd := exec.Command("docker-compose", "up", "-d", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test() {
	mg.Deps(sheazuzu.Test)
}
