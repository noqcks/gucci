package main

import (
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var funcMap = getFuncMap()

func getFuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	f["shell"] = shell

	return f
}

func shell(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err.Error()
	}
	output := strings.TrimSpace(string(out))
	return output
}
