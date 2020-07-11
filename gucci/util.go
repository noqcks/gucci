package gucci

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"
)

func isJsonFile(path string) bool {
	path = strings.ToLower(path)
	return strings.HasSuffix(path, "json")
}

func isYamlFile(path string) bool {
	path = strings.ToLower(path)
	return strings.HasSuffix(path, "yaml") ||
		strings.HasSuffix(path, "yml")
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
