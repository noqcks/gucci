package main

import (
	"encoding/base64"
	"os/exec"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"b64enc": b64enc,
	"join":   strings.Join,
	"split":  strings.Split,
	"shell":  shell,
}

func b64enc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func shell(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err.Error()
	}
	output := strings.TrimSpace(string(out))
	return output
}
