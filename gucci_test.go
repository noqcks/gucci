package main

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

var testMap = map[string]string{
	"TEST":     "green",
	"BACKENDS": "server1.com,server2.com",
}

func TestSub(t *testing.T) {
	tpl := `{{ .TEST }}`
	if err := runTest(tpl, "green"); err != nil {
		t.Error(err)
	}
}

func TestSubSplit(t *testing.T) {
	tpl := `{{- range split .BACKENDS "," }}{{ . }}{{- end }}`
	if err := runTest(tpl, "server1.comserver2.com"); err != nil {
		t.Error(err)
	}
}

func runTest(tpl, expect string) error {

	t, err := template.New("test").Funcs(funcMap).Parse(tpl)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	err = t.Execute(&b, testMap)
	if err != nil {
		return err
	}
	if b.String() != expect {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
