package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

var funcMap = template.FuncMap{
	"split": strings.Split,
	"shell": shell,
}

func shell(cmd string) string {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		return err.Error()
	}
	output := strings.TrimSpace(string(out))
	return output
}

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

func noArgs() bool {
	if len(os.Args) < 2 {
		return true
	}
	return false
}

func ExecuteTemplate(valuesIn map[string]string, out io.Writer, tplFile string) error {
	tplName := filepath.Base(tplFile)
	tpl, err := template.New(tplName).Funcs(funcMap).ParseFiles(tplFile)
	if err != nil {
		return fmt.Errorf("Error parsing template(s): %v", err)
	}
	err = tpl.Execute(out, valuesIn)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "gucci"
	app.Usage = "simple golang cli templating"
	app.Version = "0.0.2"
	app.Action = func(c *cli.Context) error {
		if noArgs() {
			cli.ShowAppHelp(c)
			return nil
		}
		err := ExecuteTemplate(Env(), os.Stdout, os.Args[1])
		if err != nil {
			return err
		}
		return nil
	}
	app.Run(os.Args)
}
