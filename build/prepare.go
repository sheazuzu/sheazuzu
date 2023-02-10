/*
 *  prepare.go
 *  Created on 08.11.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package build

import (
	"github.com/pkg/errors"
	"html/template"
	"os"
	"path"
)

type VersionData struct {
	Version string
}

func PrepareVersion(version, inputFile, outputFilePath string) error {

	versionData := VersionData{version}

	tpl, err := template.New(path.Base(inputFile)).ParseFiles(inputFile)
	if err != nil {
		return errors.Wrap(err, "error reading template")
	}
	f, err := os.Create(outputFilePath)
	if err != nil {
		return errors.Wrap(err, "error creating output file")
	}

	err = tpl.Execute(f, versionData)
	if err != nil {
		return errors.Wrap(err, "error executing template")
	}

	return nil
}
