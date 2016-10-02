package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

func getkeyval(item string) (key, val string) {
	splits := strings.Split(item, "=")
	key = splits[0]
	val = strings.Join(splits[1:], "=")
	return key, val
}

func Env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		key, val := getkeyval(i)
		env[key] = val
	}
	return env
}

func ExecuteTemplates(values_in map[string]string, out io.Writer, tpl_file string) error {
	funcMap := template.FuncMap{
		"split": strings.Split,
	}

	tpl, err := template.New(tpl_file).Funcs(funcMap).ParseFiles(tpl_file)
	if err != nil {
		return fmt.Errorf("Error parsing template(s): %v", err)
	}

	err = tpl.Execute(out, values_in)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func main() {
	err := ExecuteTemplates(Env(), os.Stdout, os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
