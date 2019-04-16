package main

import (
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
)

var funcMap = getFuncMap()

func getFuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	f["shell"] = shell
	f["toYaml"] = toYaml

	return f
}

func shell(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic("Issue running command: " + err.Error())
		return ""
	}
	output := strings.TrimSpace(string(out))
	return output
}

func toYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic("Issue marshalling yaml: " + err.Error())
		return ""
	}
	return string(data)
}
