package sheazuzu

import (
	"log"
	"sheazuzu/build"
)

var (
	MODULE  = "sheazuzu"
	VERSION string
)

func init() {
	props, err := build.ReadProperties(build.GetModuleDir(MODULE) + "/../versions.properties")
	if err != nil {
		log.Fatalf("Failed to read the versions.properties file: %s", err.Error())
	}
	VERSION = props["sheazuzu.version"]
}

func Prepare() error {
	return build.PrepareVersion(VERSION, build.GetAPIDir(MODULE)+"/swagger.yaml", build.GetTargetDir(MODULE)+"/swagger-sheazuzu.yaml")
}

func GenerateServer() error {
	return build.GenerateSwaggerServer(build.GetTargetDir(MODULE)+"/swagger-sheazuzu.yaml", "sheazuzu", build.GetGeneratedDir(MODULE)+"/sheazuzu")
}

func Build() error {
	return build.Build(MODULE, "sheazuzu", VERSION, build.LINUX, build.WINDOWS, build.MAC)
}

func Test() error {
	return build.Test(MODULE)
}

func TestCI() error {
	return build.TestCI(MODULE)
}

func Clean() error {
	return build.Clean(MODULE)
}
