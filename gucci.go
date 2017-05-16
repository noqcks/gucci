package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"log"

	"github.com/urfave/cli"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	app := cli.NewApp()
	app.Name = "gucci"
	app.Usage = "simple CLI templating"
	app.Version = "0.0.4"
	app.Action = func(c *cli.Context) error {
		f := ""
		if !noArgs() {
			f = os.Args[1]
		}
		return run(f)
	}
	app.Run(os.Args)
}


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
		logError("Error occurred while loading data:", err)
	}

	err = executeTemplate(env(), os.Stdout, tpl)
	if err != nil {
		logError("Error occurred while attempting to template:", err)
	}
	return nil
}

func logError(msg string, err error) {
  logger.Println(msg)
  logger.Println(err.Error())
}
