package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func getFuncMap(t *template.Template) template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	f["include"] = include(t)
	f["shell"] = shell
	f["toYaml"] = toYaml

	return f
}

func include(t *template.Template) func(templateName string, vars ...interface{}) (string, error) {
	return func(templateName string, vars ...interface{}) (string, error) {
		if len(vars) > 1 {
			return "", errors.New(fmt.Sprintf("Call to include may pass zero or one vars structure, got %v.", len(vars)))
		}
		buf := bytes.NewBuffer(nil)
		included := t.Lookup(templateName)
		if included == nil {
			return "", errors.New(fmt.Sprintf("No such template '%v' found while calling 'include'.", templateName))
		}

		if err := included.ExecuteTemplate(buf, templateName, vars[0]); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

func shell(cmd ...string) (string, error) {
	out, err := exec.Command("bash", "-c", strings.Join(cmd[:], "")).Output()
	output := strings.TrimSpace(string(out))
	if err != nil {
		return "", errors.Wrap(err, "Issue running command: "+output)
	}

	return output, nil
}

func toYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", errors.Wrap(err, "Issue marsahling yaml")
	}
	return string(data), nil
}
