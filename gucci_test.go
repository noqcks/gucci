package main

import (
	"bytes"
	"fmt"
	"reflect"
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
	err = executeTemplate(testVarMap, &b, tpl, []string{})
	if err != nil {
		return err
	}
	if b.String() != expect {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}

func Test_getTplOpt(t *testing.T) {
	type args struct {
		cliTplOpt []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Default value if no option specified",
			args: args {
				cliTplOpt: []string{},
			},
			want: []string{"missingkey=error"},
		},
		{
			name: "Specified Option",
			args: args {
				cliTplOpt: []string{"missingkey=zero"},
			},
			want: []string{"missingkey=zero"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTplOpt(tt.args.cliTplOpt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTplOpt() = %v, want %v", got, tt.want)
			}
	})
}
}
