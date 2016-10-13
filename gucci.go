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
	"io/ioutil"
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

func LoadFile(tplFile string) (*template.Template, error) {
	tplName := filepath.Base(tplFile)
	tpl, err := template.New(tplName).Funcs(funcMap).ParseFiles(tplFile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

func LoadStream(name string, in io.Reader) (*template.Template, error) {
	tplBytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading template(s): %v", err)
	}
	tplStr := string(tplBytes)
	return LoadString(name, tplStr)
}

func LoadString(name, s string) (*template.Template, error) {
	tpl, err := template.New(name, ).Funcs(funcMap).Parse(s)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

func LoadFileOrStdin(f string) (*template.Template, error) {
	var tpl *template.Template
	if f == "" {
		t, err := LoadStream("-", os.Stdin)
		if err != nil {
			return nil, err
		}
		tpl = t
	} else {
		t, err := LoadFile(f)
		if err != nil {
			return nil, err
		}
		tpl = t
	}
	return tpl, nil
}

func ExecuteTemplate(valuesIn map[string]string, out io.Writer, tpl *template.Template) error {
	err := tpl.Execute(out, valuesIn)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func Run(arg string) error {
	tpl, err := LoadFileOrStdin(arg)
	if err != nil {
		return err
	}

	err = ExecuteTemplate(Env(), os.Stdout, tpl)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "gucci"
	app.Usage = "simple golang cli templating"
	app.Version = "0.0.3"
	app.Action = func(c *cli.Context) error {
		f := ""
		if !noArgs() {
			f = os.Args[1]
		}
		return Run(f)
	}
	app.Run(os.Args)
}
