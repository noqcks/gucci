package gucci

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Idea is to collect public methods here

// Env ...
func Env() map[string]interface{} {
	env := make(map[string]interface{})
	for _, i := range os.Environ() {
		key, val := GetKeyVal(i)
		env[key] = val
	}
	return env
}

// ExecuteTemplate ...
func ExecuteTemplate(valuesIn map[string]interface{}, out io.Writer, tpl *template.Template) error {
	tpl.Option("missingkey=error")
	err := tpl.Execute(out, valuesIn)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

// GetKeyVal ...
func GetKeyVal(item string) (key, val string) {
	splits := strings.Split(item, "=")
	key = splits[0]
	val = strings.Join(splits[1:], "=")
	return key, val
}

// KeyValToMap ...
func KeyValToMap(key, val string) map[string]interface{} {
	parts := strings.Split(key, ".")

	// Reverse order
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	m := map[string]interface{}{
		parts[0]: val,
	}

	for _, part := range parts[1:] {
		m = map[string]interface{}{
			part: m,
		}
	}

	return m
}

// LoadTemplateFile ...
func LoadTemplateFile(tplFile string) (*template.Template, error) {
	tplName := filepath.Base(tplFile)
	tpl, err := template.New(tplName).Funcs(funcMap).ParseFiles(tplFile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

// LoadTemplateFileOrStdin ...
func LoadTemplateFileOrStdin(f string) (*template.Template, error) {
	var tpl *template.Template
	if f == "" {
		t, err := LoadTemplateStream("-", os.Stdin)
		if err != nil {
			return nil, err
		}
		tpl = t
	} else {
		t, err := LoadTemplateFile(f)
		if err != nil {
			return nil, err
		}
		tpl = t
	}
	return tpl, nil
}

// LoadTemplateStream ...
func LoadTemplateStream(name string, in io.Reader) (*template.Template, error) {
	tplBytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading template(s): %v", err)
	}
	tplStr := string(tplBytes)
	return LoadTemplateString(name, tplStr)
}

// LoadTemplateString ...
func LoadTemplateString(name, s string) (*template.Template, error) {
	tpl, err := template.New(name).Funcs(funcMap).Parse(s)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

// LoadVarsFile ...
func LoadVarsFile(path string) (map[string]interface{}, error) {
	var result map[string]interface{}
	var err error

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if isJsonFile(path) {
		result, err = unmarshalJsonFile(content)
	} else if isYamlFile(path) {
		result, err = unmarshalYamlFile(content)
	} else {
		err = fmt.Errorf("unsupported variables file type: %s", path)
	}

	if err != nil {
		return nil, err
	}

	return result, err
}
