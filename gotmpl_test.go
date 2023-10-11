package main

import (
	"os"
	"testing"
)

func TestGotmpl(t *testing.T) {
	valuesJson := `{
	"hello": "world",
	"number": {
		"num": 5,
		"string": "five"
	},
	"cmd": "ls;echo gotcha"
}`
	testTemplate := "Hello {{ .hello }}. Number {{.number.num}} is {{.number.string}}. Command is {{ .cmd | shellescape }}.\n"

	v, err := os.CreateTemp("", "gotmpl-values-*.json")
	if err != nil {
		t.Fatal(err)
	}
	vname := v.Name()
	defer os.Remove(vname)
	if _, err := v.WriteString(valuesJson); err != nil {
		v.Close()
		t.Fatal(err)
	}
	v.Close()

	tmpl, err := os.CreateTemp("", "gotmpl-file-*.template")
	if err != nil {
		t.Fatal(err)
	}
	tmplname := tmpl.Name()
	defer os.Remove(tmplname)
	if _, err := tmpl.WriteString(testTemplate); err != nil {
		tmpl.Close()
		t.Fatal(err)
	}
	tmpl.Close()

	out, err := os.CreateTemp("", "gotmpl-output-*")
	if err != nil {
		t.Fatal(err)
	}
	outname := out.Name()
	out.Close()
	defer os.Remove(outname)

	if err := gotmpl(vname, tmplname, outname); err != nil {
		t.Fatal(err)
	}

	result, err := os.ReadFile(outname)
	if err != nil {
		t.Fatal(err)
	}

	if got, expected := string(result), "Hello world. Number 5 is five. Command is 'ls;echo gotcha'.\n"; got != expected {
		t.Fatalf("Expected %q, but got %q", expected, got)
	}

}
