package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		key, val := getkeyval(i)
		env[key] = val
	}
	return env
}

func getkeyval(item string) (key, val string) {
	splits := strings.Split(item, "=")
	key = splits[0]
	val = strings.Join(splits[1:], "=")
	return key, val
}

func noArgs() bool {
	if len(os.Args) < 2 {
		return true
	}
	return false
}

func loadFile(tplFile string) (*template.Template, error) {
	tplName := filepath.Base(tplFile)
	tpl, err := template.New(tplName).Funcs(funcMap).ParseFiles(tplFile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

func loadStream(name string, in io.Reader) (*template.Template, error) {
	tplBytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading template(s): %v", err)
	}
	tplStr := string(tplBytes)
	return loadString(name, tplStr)
}

func loadString(name, s string) (*template.Template, error) {
	tpl, err := template.New(name).Funcs(funcMap).Parse(s)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

func loadFileOrStdin(f string) (*template.Template, error) {
	var tpl *template.Template
	if f == "" {
		t, err := loadStream("-", os.Stdin)
		if err != nil {
			return nil, err
		}
		tpl = t
	} else {
		t, err := loadFile(f)
		if err != nil {
			return nil, err
		}
		tpl = t
	}
	return tpl, nil
}
