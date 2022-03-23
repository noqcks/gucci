package main

import (
	"bytes"
	"fmt"
	"strings"
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

func TestFuncIncludeNoVars(t *testing.T) {
	tpl := `{{ define "a" }}Hi Jane{{ end }}{{ include "a" . }}`
	if err := runTest(tpl, "Hi Jane"); err != nil {
		t.Error(err)
	}
}

func TestFuncIncludeWithVars(t *testing.T) {
	tpl := `{{ define "a" }}Hi {{ .name }}{{ end }}{{ $_ := set . "name" "John" }}{{ include "a" . }}`
	if err := runTest(tpl, "Hi John"); err != nil {
		t.Error(err)
	}
}

func TestFuncIncludePipe(t *testing.T) {
	tpl := `{{ define "a" }}Hi Jane{{ end }}{{ include "a" . | indent 2 }}`
	if err := runTest(tpl, "  Hi Jane"); err != nil {
		t.Error(err)
	}
}

func TestFuncIncludeTooManyVars(t *testing.T) {
	tpl := `{{ define "a" }}Hi {{ .name }}{{ end }}{{ $_ := set . "name" "John" }}{{ include "a" . . }}`
	if err := runTest(tpl, ""); err == nil {
		t.Error("expected error to many vars")
	}
}

func TestFuncIncludeBadName(t *testing.T) {
	tpl := `{{ define "a" }}Hi Jane{{ end }}{{ include "b" . }}`
	if err := runTest(tpl, ""); err == nil {
		t.Error("expected error bad template name")
	}
}

func TestFuncShell(t *testing.T) {
	tpl := `{{ shell "echo hello" }}`
	if err := runTest(tpl, "hello"); err != nil {
		t.Error(err)
	}
}

func TestFuncShellArguments(t *testing.T) {
	tpl := `{{ shell "echo " "hello" "world"}}`
	if err := runTest(tpl, "helloworld"); err != nil {
		t.Error(err)
	}
}

func TestFuncShellError(t *testing.T) {
	tpl := `{{ shell "non-existent" }}`
	if err := runTest(tpl, ""); err == nil {
		t.Error("expected error missing")
	}
}

func TestFuncShellDetailedError(t *testing.T) {
	tpl := `{{ shell "echo saboteur ; exit 1" }}`
	err := runTest(tpl, "")
	if !strings.Contains(err.Error(), "saboteur") {
		t.Error("expected stdout in error missing. actual: ", err)
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
	err = executeTemplate(testVarMap, &b, tpl, []string{})
	if err != nil {
		return err
	}
	if b.String() != expect {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
