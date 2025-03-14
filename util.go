package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func env() map[string]interface{} {
	env := make(map[string]interface{})
	for _, i := range os.Environ() {
		key, val := getKeyVal(i)
		env[key] = val
	}
	return env
}

func getKeyVal(item string) (key, val string) {
	splits := strings.Split(item, "=")
	key = splits[0]
	val = strings.Join(splits[1:], "=")
	return key, val
}

func loadTemplateFile(tplFile string) (*template.Template, error) {
	tplName := filepath.Base(tplFile)
	tpl := template.New(tplName)
	_, err := tpl.Funcs(getFuncMap(tpl)).ParseFiles(tplFile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

func loadTemplateStream(name string, in io.Reader) (*template.Template, error) {
	tplBytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading template(s): %v", err)
	}
	tplStr := string(tplBytes)
	return loadTemplateString(name, tplStr)
}

func loadTemplateString(name, s string) (*template.Template, error) {
	tpl := template.New(name)
	tpl, err := tpl.Funcs(getFuncMap(tpl)).Parse(s)
	if err != nil {
		return nil, fmt.Errorf("Error parsing template(s): %v", err)
	}
	return tpl, nil
}

func loadTemplateFileOrStdin(f string) (*template.Template, error) {
	var tpl *template.Template
	if f == "" {
		t, err := loadTemplateStream("-", os.Stdin)
		if err != nil {
			return nil, err
		}
		tpl = t
	} else {
		t, err := loadTemplateFile(f)
		if err != nil {
			return nil, err
		}
		tpl = t
	}
	return tpl, nil
}

func isJsonFile(path string) bool {
	path = strings.ToLower(path)
	return strings.HasSuffix(path, "json")
}

func isYamlFile(path string) bool {
	path = strings.ToLower(path)
	return strings.HasSuffix(path, "yaml") ||
		strings.HasSuffix(path, "yml")
}

func loadVarsFile(path string) (map[string]interface{}, error) {
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

func unmarshalJsonFile(content []byte) (map[string]interface{}, error) {
	var vars map[string]interface{}
	err := json.Unmarshal(content, &vars)
	if err != nil {
		return nil, err
	}
	return vars, nil
}

func unmarshalYamlFile(content []byte) (map[string]interface{}, error) {
	var vars map[string]interface{}
	err := yaml.Unmarshal(content, &vars)
	if err != nil {
		return nil, err
	}
	return vars, nil
}

func keyValToMap(key, val string) map[string]interface{} {
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
