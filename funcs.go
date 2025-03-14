package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func getFuncMap(t *template.Template) template.FuncMap {
	f := sprig.TxtFuncMap()

	f["include"] = include(t)
	f["shell"] = shell
	f["toYaml"] = toYaml

	f["toJson"] = toJson
	f["mustToJson"] = mustToJson

	return f
}

// toJson converts a value to JSON with special handling for map[interface{}]interface{}
// This function overrides Sprig's toJson function to properly handle YAML-parsed data.
// When working with templates, we often load data from YAML files which creates maps with
// interface{} keys, but JSON only supports string keys. This implementation ensures that
// data loaded from YAML (or other sources) with non-string keys is properly converted
// before JSON serialization, preventing common "json: unsupported type" errors when
// templates try to convert complex nested structures to JSON.
func toJson(v interface{}) (string, error) {
	// Convert any map[interface{}]interface{} to map[string]interface{}
	jsonCompatible := convertToJSONCompatible(v)
	data, err := json.Marshal(jsonCompatible)
	if err != nil {
		return "", errors.Wrap(err, "error calling toJson")
	}
	return string(data), nil
}

// mustToJson is like toJson but panics on error
// we need to override this because sprig's mustToJson doesn't handle
// map[interface{}]interface{} which is what we get from yaml
func mustToJson(v interface{}) string {
	s, err := toJson(v)
	if err != nil {
		panic(err)
	}
	return s
}

// convertToJSONCompatible converts YAML parsed data (with interface{} keys)
// to data with only string keys for JSON compatibility
func convertToJSONCompatible(v interface{}) interface{} {
	switch v := v.(type) {
	case map[interface{}]interface{}:
		// Convert map with interface{} keys to map with string keys
		result := make(map[string]interface{})
		for k, v := range v {
			result[fmt.Sprintf("%v", k)] = convertToJSONCompatible(v)
		}
		return result
	case []interface{}:
		// Convert each item in the slice
		for i, item := range v {
			v[i] = convertToJSONCompatible(item)
		}
	}
	return v
}

func include(t *template.Template) func(templateName string, vars ...interface{}) (string, error) {
	return func(templateName string, vars ...interface{}) (string, error) {
		if len(vars) > 1 {
			return "", errors.New(fmt.Sprintf("Call to include may pass zero or one vars structure, got %v.", len(vars)))
		}
		buf := bytes.NewBuffer(nil)
		included := t.Lookup(templateName)
		if included == nil {
			return "", errors.New(fmt.Sprintf("No such template '%v' found while calling 'include'.", templateName))
		}

		if err := included.ExecuteTemplate(buf, templateName, vars[0]); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

func shell(cmd ...string) (string, error) {
	out, err := exec.Command("bash", "-c", strings.Join(cmd[:], "")).Output()
	output := strings.TrimSpace(string(out))
	if err != nil {
		return "", errors.Wrap(err, "Issue running command: "+output)
	}

	return output, nil
}

func toYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", errors.Wrap(err, "Issue marsahling yaml")
	}
	return string(data), nil
}
