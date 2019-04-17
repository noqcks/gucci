package main

import (
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
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

func shell(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", errors.Wrap(err, "Issue running command")
	}
	output := strings.TrimSpace(string(out))
	return output, nil
}

func toYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", errors.Wrap(err, "Issue marsahling yaml")
	}
	return string(data), nil
}
