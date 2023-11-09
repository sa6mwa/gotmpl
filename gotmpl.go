package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/alessio/shellescape"
)

var output string

var additionalFunctions template.FuncMap = template.FuncMap{
	"shellescape": func(cmd string) string {
		return shellescape.Quote(cmd)
	},
}

func tmplFuncSignature(key string, fn any) string {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return ""
	}
	signature := key + " "
	for i := 0; i < fnType.NumIn(); i++ {
		if i > 0 {
			signature += ", "
		}
		argType := fnType.In(i)
		signature += argType.String()
	}
	signature += " (returns "
	for i := 0; i < fnType.NumOut(); i++ {
		if i > 0 {
			signature += ", "
		}
		outType := fnType.Out(i)
		signature += outType.String()
	}
	signature += ")"
	return signature
}

func enumFunctions() string {
	out := []string{}
	for k, fn := range additionalFunctions {
		out = append(out, tmplFuncSignature(k, fn))
	}
	return strings.Join(out, ", ")
}

func gotmpl(valuesJson, fileToTemplate, outputFile string) error {
	valuesF := os.Stdin
	if valuesJson != "" {
		var err error
		valuesF, err = os.Open(valuesJson)
		if err != nil {
			return err
		}
		defer valuesF.Close()
	}

	templateF := os.Stdin
	if fileToTemplate != "" {
		var err error
		templateF, err = os.Open(fileToTemplate)
		if err != nil {
			return err
		}
		defer templateF.Close()
	}

	if valuesF == os.Stdin && templateF == os.Stdin {
		return errors.New("values.json and template file can not both be stdin")
	}

	var values map[string]interface{}

	if err := json.NewDecoder(valuesF).Decode(&values); err != nil {
		return err
	}

	outputF := os.Stdout
	if outputFile != "" {
		var err error
		outputF, err = os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer outputF.Close()
	}

	templateData, err := io.ReadAll(templateF)
	if err != nil {
		return err
	}

	functions := sprig.TxtFuncMap()
	for k, fnc := range additionalFunctions {
		functions[k] = fnc
	}

	tmpl, err := template.New("template").Funcs(functions).Parse(string(templateData))
	if err != nil {
		return err
	}
	return tmpl.Execute(outputF, &values)
}

func main() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "[option] [values_json] file_to_template")
		fmt.Fprintln(os.Stderr, "If values_json is omitted, values are read from stdin (as json)")
		fmt.Fprintln(os.Stderr, "Sprig functions are supported, see https://masterminds.github.io/sprig/")
		fmt.Fprintln(os.Stderr, os.Args[0], "also supports the following aditional function(s):")
		fmt.Fprintln(os.Stderr, enumFunctions())
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
	flag.StringVar(&output, "o", "", "output file, default is stdout")
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() < 2 {
		if err := gotmpl("", flag.Arg(0), output); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	} else {
		if err := gotmpl(flag.Arg(0), flag.Arg(1), output); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	}
}
