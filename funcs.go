package main

import (
	"os/exec"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"split": strings.Split,
	"shell": shell,
}

func shell(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err.Error()
	}
	output := strings.TrimSpace(string(out))
	return output
}
