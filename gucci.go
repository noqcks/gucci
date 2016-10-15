package main

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/urfave/cli"
)

func executeTemplate(valuesIn map[string]string, out io.Writer, tpl *template.Template) error {
	err := tpl.Execute(out, valuesIn)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func run(arg string) error {
	tpl, err := loadFileOrStdin(arg)
	if err != nil {
		return err
	}

	err = executeTemplate(env(), os.Stdout, tpl)
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
		return run(f)
	}
	app.Run(os.Args)
}
