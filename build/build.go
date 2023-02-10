package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type OS int

const (
	LINUX OS = iota
	MAC
	WINDOWS
)

var goos = map[OS]string{
	LINUX:   "linux",
	MAC:     "darwin",
	WINDOWS: "windows",
}

func Build(module, outputName, version string, operatingSystems ...OS) error {
	if len(operatingSystems) == 0 {
		operatingSystems = []OS{LINUX}
	}

	for _, ops := range operatingSystems {
		fmt.Printf("Building %s %s for %s\n", module, version, goos[ops])

		filename := fmt.Sprintf("%s-%s-%s", outputName, goos[ops], "amd64")
		if ops == WINDOWS {
			filename += ".exe"
		}

		cmd := exec.Command("go", "build", "-o", "../target/"+filename, "-ldflags", "-X main.ProjectVersion="+version, ".")
		cmd.Dir = GetSourceDir(module)

		cmd.Env = append(cmd.Env, os.Environ()...)
		cmd.Env = append(cmd.Env, "GOOS="+goos[ops])
		cmd.Env = append(cmd.Env, "GOARCH=amd64")
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")

		out := bytes.NewBuffer([]byte{})
		cmd.Stdout = out
		cmd.Stderr = out

		err := cmd.Run()

		fmt.Fprint(os.Stdout, out.String())

		if err != nil {
			return err
		}
	}

	return nil
}

func Test(module string) error {
	fmt.Printf("Testing %s...\n", module)

	_ = GetTargetDir(module)

	cmd := exec.Command("go", "test", "./...", "--coverprofile=../target/cover.out")
	cmd.Dir = GetSourceDir(module)

	out := bytes.NewBuffer([]byte{})
	cmd.Stdout = out
	cmd.Stderr = out

	cmdError := cmd.Run()
	fmt.Fprint(os.Stdout, out.String())

	return cmdError
}

func TestCI(module string) error {
	fmt.Printf("Testing %s...\n", module)

	_ = GetTargetDir(module)

	cmd := exec.Command("go", "test", "./...", "--coverprofile=../target/cover.out", "-json")
	cmd.Dir = GetSourceDir(module)

	out := bytes.NewBuffer([]byte{})
	cmd.Stdout = out
	cmd.Stderr = out

	cmdError := cmd.Run()
	fmt.Fprint(os.Stdout, out.String())

	err := ioutil.WriteFile(GetTargetDir(module)+"/report.out", out.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return cmdError
}

func Clean(module string) error {
	fmt.Printf("Cleaning %s...\n", module)
	err := os.RemoveAll(GetTargetDir(module))
	if err != nil {
		return err
	}
	return os.RemoveAll(GetGeneratedDir(module))
}
