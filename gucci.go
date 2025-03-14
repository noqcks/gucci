package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/imdario/mergo"

	"github.com/urfave/cli"
)

var logger = log.New(os.Stderr, "", 0)

const (
	flagSetVar     = "s"
	flagSetVarLong = flagSetVar + ",set-var"

	flagVarsFile     = "f"
	flagVarsFileLong = flagVarsFile + ",vars-file"

	flagSetOpt = "o"
	flagSetOptLong = flagSetOpt + ",tpl-opt"
)

var (
	AppVersion = "0.0.0-dev.0" // Injected
)

func main() {
	app := cli.NewApp()
	app.Name = "gucci"
	app.Usage = "simple CLI Go lang templating"
	app.UsageText = app.Name + " [options] [template]"
	app.Version = AppVersion

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  flagSetVarLong,
			Usage: "A `KEY=VALUE` pair variable",
		},
		cli.StringSliceFlag{
			Name:  flagVarsFileLong,
			Usage: "A json or yaml `FILE` from which to read variables (can be specified multiple times)",
		},
		cli.StringSliceFlag{
			Name: flagSetOptLong,
			Usage: "A template option (`KEY=VALUE`) to be applied",
			Value: &cli.StringSlice{"missingkey=error"},
		},
	}

	app.Action = func(c *cli.Context) error {
		tplPath := c.Args().First()
		vars, err := loadVariables(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		err = run(tplPath, vars, c.StringSlice(flagSetOpt))
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	}
	app.Run(os.Args)
}

func loadInputVarsFile(c *cli.Context) (map[string]interface{}, error) {
	vars := make(map[string]interface{})

	varsFiles := c.StringSlice(flagVarsFile)
	for _, varsFilePath := range varsFiles {
		if varsFilePath != "" {
			v, err := loadVarsFile(varsFilePath)
			if err != nil {
				return nil, err
			}

			err = mergo.Merge(&vars, v, mergo.WithOverride)
			if err != nil {
				return nil, err
			}
		}
	}

	return vars, nil
}

func loadInputVarsOptions(c *cli.Context) (map[string]interface{}, error) {

	vars := make(map[string]interface{})

	for _, varStr := range c.StringSlice(flagSetVar) {
		key, val := getKeyVal(varStr)
		varMap := keyValToMap(key, val)

		err := mergo.Merge(&vars, varMap, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	return vars, nil
}

func loadVariables(c *cli.Context) (map[string]interface{}, error) {

	vars, err := loadInputVarsFile(c)
	if err != nil {
		return nil, err
	}

	envVars := env()
	err = mergo.Merge(&vars, envVars, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	optVars, err := loadInputVarsOptions(c)
	if err != nil {
		return nil, err
	}

	err = mergo.Merge(&vars, optVars, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

func executeTemplate(valuesIn map[string]interface{}, out io.Writer, tpl *template.Template, opt []string) error {
	tpl.Option(opt...)
	err := tpl.Execute(out, valuesIn)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func run(tplPath string, vars map[string]interface{}, tplOpt []string) error {
	tpl, err := loadTemplateFileOrStdin(tplPath)
	if err != nil {
		return err
	}

	err = executeTemplate(vars, os.Stdout, tpl, tplOpt)
	if err != nil {
		return err
	}

	return nil
}

func logError(msg string, err error) {
	logger.Println(msg)
	logger.Println(err.Error())
}
