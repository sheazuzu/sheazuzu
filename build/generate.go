/*
 *  generate.go
 *  Created on 08.11.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package build

import (
	"bytes"
	"fmt"
	"github.com/magefile/mage/target"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/codegen"
	"github.com/deepmap/oapi-codegen/pkg/util"
	"github.com/pkg/errors"
)

func ConvertToOpenApi3(inputFilePath, outputFilePath string) error {

	// only convert if the swagger file has changed
	modified, err := target.Path(outputFilePath, inputFilePath)
	if err != nil {
		return err
	}

	if !modified {
		return nil
	}

	fmt.Printf("Converting swagger file '%s' to OpenAPI 3.0 ... \n", inputFilePath)

	input, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"https://converter.swagger.io/api/convert",
		"application/yaml",
		bytes.NewBuffer(input),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFilePath, output, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func GenerateSwaggerClient(api, packageName, outputFilePath string) error {
	opts := codegen.Options{
		GenerateClient: true,
		GenerateTypes:  true,
	}

	swagger, err := util.LoadSwagger(api)
	if err != nil {
		return errors.Wrap(err, "could not load api")
	}

	code, err := codegen.Generate(swagger, packageName, opts)
	if err != nil {
		return errors.Wrap(err, "error generating code")
	}

	_ = os.MkdirAll(outputFilePath, os.ModePerm)

	err = ioutil.WriteFile(outputFilePath+"/client.go", []byte(code), 0644)
	if err != nil {
		return errors.Wrap(err, "error writing generating code")
	}

	return nil
}

func GenerateSwaggerServer(api, packageName, outputFilePath string) error {
	opts := codegen.Options{
		GenerateChiServer: true,
		EmbedSpec:         true,
		GenerateTypes:     true,
		UserTemplates: map[string]string{
			"chi-interface.tmpl": chiInterface,
			"chi-handler.tmpl":   chiHandler,
		},
	}

	swagger, err := util.LoadSwagger(api)
	if err != nil {
		return errors.Wrap(err, "could not load api")
	}

	code, err := codegen.Generate(swagger, packageName, opts)
	if err != nil {
		return errors.Wrap(err, "error generating code")
	}

	_ = os.MkdirAll(outputFilePath, os.ModePerm)

	err = ioutil.WriteFile(outputFilePath+"/server.go", []byte(code), 0644)
	if err != nil {
		return errors.Wrap(err, "error writing generating code")
	}

	return nil
}

func GenerateStringFromFile(filePath, outputFilePath, stringName string) error {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "could not read file")
	}

	bs := string(b)
	bs = strings.ReplaceAll(bs, "`", "\\`")

	_ = os.MkdirAll(outputFilePath, os.ModePerm)

	w, err := os.Create(outputFilePath + "/" + stringName + ".go")
	if err != nil {
		return errors.Wrap(err, "error writing generating code")
	}

	fmt.Fprintf(w, "package generated\n\nconst %s = `%s`", stringName, bs)
	return nil
}

var chiInterface = `
	// ServerInterface represents all server handlers.
	type ServerInterface interface {
	{{range .}}{{.SummaryAsComment }}
	// ({{.Method}} {{.Path}})
	{{.OperationId}}(w http.ResponseWriter, r *http.Request{{genParamArgs .PathParams}}{{if .RequiresParamObject}}, params {{.OperationId}}Params{{end}})
	{{end}}
	}

	type ServerWithMiddleware struct {
	ServerInterface
	{{range .}}// {{.Summary | stripNewLines }} ({{.Method}} {{.Path}})
	{{.OperationId}}Middlewares chi.Middlewares
	{{end}}
	}

	func NewServerWithMiddleware(si ServerInterface) ServerWithMiddleware {
		return ServerWithMiddleware{
			ServerInterface: si,
		}
	}
`

var chiHandler = `
// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerWithMiddleware) http.Handler {
  return HandlerFromMux(si, chi.NewRouter())
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerWithMiddleware, r chi.Router) http.Handler {
{{if .}}wrapper := ServerInterfaceWrapper{
        Handler: si,
    }
{{end}}
{{range .}}r.Group(func(r chi.Router) {
  r.With(si.{{.OperationId}}Middlewares...).{{.Method | lower | title }}("{{.Path | swaggerUriToChiUri}}", wrapper.{{.OperationId}})
})
{{end}}
  return r
}
`
