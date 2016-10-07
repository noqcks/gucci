package main

import (
	"bytes"
	"fmt"
	"os"
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
	tpl := `{{ range split .BACKENDS "," }}{{ . }}{{ end }}`
	if err := runTest(tpl, "server1.comserver2.com"); err != nil {
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

func TestGetKeyVal(t *testing.T) {
	tests := []struct {
		in, k, v string
	}{
		{"k=v", "k", "v"},
		{"kv", "kv", ""},
		{"=kv", "", "kv"},
	}
	for _, tt := range tests {
		k, v := getkeyval(tt.in)
		if k != tt.k || v != tt.v {
			t.Errorf("broken behavior. Expected: %#v Got: %v %v", tt, k, v)
		}
	}
}

func TestEnv(t *testing.T) {
	os.Setenv("k", "v")
	envs := Env()
	if v, ok := envs["k"]; !ok || (ok && v != "v") {
		t.Errorf("broken behavior. Expected: %v. Got: %v", "v", v)
	}
}

func TestNoArgs(t *testing.T) {
	oldArgs := os.Args

	cases := []struct {
		args     []string
		expected bool
	}{
		{[]string{"a"}, true},
		{[]string{"a", "b", "c"}, false},
	}

	for _, tt := range cases {
		os.Args = tt.args
		if got := noArgs(); got != tt.expected {
			t.Errorf("%#v Expected: %v Got: %v", tt.args, tt.expected, got)
		}
	}
	os.Args = oldArgs
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
