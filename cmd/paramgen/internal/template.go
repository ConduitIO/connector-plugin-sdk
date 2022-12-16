package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

const (
	tmpl = `// Code generated by ParamGen. DO NOT EDIT.
// Source: github.com/conduitio/conduit-connector-sdk/cmd/paramgen

package {{ $.Package }}

import (
	{{ if $.HasRegex }}
	"regexp"
	{{ end }}
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func ({{ $.Struct }}) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
		{{- range $name, $parameter := .Parameters }}
		"{{ $name }}": {
			Default:     "{{ .Default }}",
			Description: "{{ .Description }}",
			Type:        sdk.{{ .GetTypeConstant }},
			Validations: []sdk.Validation{
				{{- range $index, $validation := .Validations }}
					{{ $parameter.GetValidation $index }},
				{{- end }}
			},
		},
		{{- end }}
	}
}
`
)

type templateData struct {
	Parameters map[string]parameter
	Package    string
	Struct     string
}

var parameterTypeConstantMapping = map[sdk.ParameterType]string{
	sdk.ParameterTypeString:   "ParameterTypeString",
	sdk.ParameterTypeInt:      "ParameterTypeInt",
	sdk.ParameterTypeFloat:    "ParameterTypeFloat",
	sdk.ParameterTypeBool:     "ParameterTypeBool",
	sdk.ParameterTypeFile:     "ParameterTypeFile",
	sdk.ParameterTypeDuration: "ParameterTypeDuration",
}

type parameter sdk.Parameter

func (p parameter) GetTypeConstant() string {
	return parameterTypeConstantMapping[p.Type]
}
func (p parameter) GetValidation(index int) string {
	validation := p.Validations[index]

	regexValidation, ok := validation.(sdk.ValidationRegex)
	if !ok {
		// default behavior
		return fmt.Sprintf("%#v", p.Validations[index])
	}

	validationType := reflect.TypeOf(validation).String()
	validationParameters := fmt.Sprintf("Regex: regexp.MustCompile(%q)", regexValidation.Regex)
	return fmt.Sprintf("%s{%s}", validationType, validationParameters)
}

func (t templateData) HasRegex() bool {
	for _, p := range t.Parameters {
		for _, v := range p.Validations {
			if _, ok := v.(sdk.ValidationRegex); ok {
				return true
			}
		}
	}
	return false
}

func GenerateCode(parameters map[string]sdk.Parameter, packageName string, structName string) string {
	// create the go template
	t := template.Must(template.New("").Parse(tmpl))

	internalParams := make(map[string]parameter, len(parameters))
	for k, v := range parameters {
		internalParams[k] = parameter(v)
	}

	data := templateData{
		Package:    packageName,
		Struct:     structName,
		Parameters: internalParams,
	}
	var processed bytes.Buffer
	// execute the template
	err := t.Execute(&processed, data)
	if err != nil {
		log.Fatalf("error executing template: %v\n", err)
	}

	// format the output as Go code in the “gofmt” style
	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Could not format processed template: %v\n", err)
	}

	return string(formatted)
}

func WriteCodeToFile(str string, path string, output string) {
	// create output directory if it does not exist
	outputDir := strings.TrimSuffix(path, "/")

	// create the output file and write data
	outputPath := outputDir + "/" + output
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("could not create file: %v\n", err)
	}
	w := bufio.NewWriter(f)
	_, err = w.WriteString(str)
	if err != nil {
		log.Fatalf("error writing to a file: %v\n", err)
	}
	w.Flush()
	log.Printf("output file created: %s\n", outputPath)
}
