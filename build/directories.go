package build

import "os"

func GetModuleDir(module string) string {
	return module
}

func GetTargetDir(module string) string {
	target := GetModuleDir(module) + "/target"

	_ = os.MkdirAll(target, os.ModePerm)

	return target
}

func GetSourceDir(module string) string {
	return GetModuleDir(module) + "/src"
}

func GetGeneratedDir(module string) string {
	generated := GetSourceDir(module) + "/generated"

	_ = os.MkdirAll(generated, os.ModePerm)

	return generated
}

func GetDocsDir(module string) string {
	return GetModuleDir(module) + "/docs"
}

func GetAPIDir(module string) string {
	return GetModuleDir(module) + "/api"
}
