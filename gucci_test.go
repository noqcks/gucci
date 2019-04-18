package main

import (
	"bytes"
	"fmt"
	"testing"
)

var testVarMap = map[string]interface{}{
	"TEST": "green",
}

func TestSub(t *testing.T) {
	tpl := `{{ .TEST }}`
	if err := runTest(tpl, "green"); err != nil {
		t.Error(err)
	}
}

func TestFuncShell(t *testing.T) {
	tpl := `{{ shell "echo hello" }}`
	if err := runTest(tpl, "hello"); err != nil {
		t.Error(err)
	}
}

func TestFuncShellError(t *testing.T) {
	tpl := `{{ shell "non-existent" }}`
	if err := runTest(tpl, ""); err == nil {
		t.Error("expected error missing")
	}
}

func TestFuncShellPipe(t *testing.T) {
	tpl := `{{ shell "echo foo | grep foo" }}`
	if err := runTest(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestFuncToYaml(t *testing.T) {
	tpl := `{{ list "a" "b" "c" | toYaml }}`
	if err := runTest(tpl, "- a\n- b\n- c\n"); err != nil {
		t.Error(err)
	}
}

func runTest(str, expect string) error {
	tpl, err := loadTemplateString("test", str)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = executeTemplate(testVarMap, &b, tpl)
	if err != nil {
		return err
	}
	if b.String() != expect {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
