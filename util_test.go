package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetKeyVal(t *testing.T) {
	tests := []struct {
		in, k, v string
	}{
		{"k=v", "k", "v"},
		{"kv", "kv", ""},
		{"=kv", "", "kv"},
	}
	for _, tt := range tests {
		k, v := getKeyVal(tt.in)
		if k != tt.k || v != tt.v {
			t.Errorf("broken behavior. Expected: %#v Got: %v %v", tt, k, v)
		}
	}
}

func TestEnv(t *testing.T) {
	os.Setenv("k", "v")
	envs := env()
	if v, ok := envs["k"]; !ok || (ok && v != "v") {
		t.Errorf("broken behavior. Expected: %v. Got: %v", "v", v)
	}
}

func TestKeyValToMap(t *testing.T) {
	tests := []struct {
		key      string
		value    string
		expected map[string]interface{}
	}{
		{
			"foo",
			"bar",
			map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			"foo1.foo2",
			"bar",
			map[string]interface{}{
				"foo1": map[string]interface{}{
					"foo2": "bar",
				},
			},
		},
		{
			"foo1.foo2.foo3",
			"bar",
			map[string]interface{}{
				"foo1": map[string]interface{}{
					"foo2": map[string]interface{}{
						"foo3": "bar",
					},
				},
			},
		},
	}
	for _, test := range tests {
		r := keyValToMap(test.key, test.value)
		if !reflect.DeepEqual(r, test.expected) {
			t.Errorf("broken behavior. Expected: %v. Got: %v", test.expected, r)
		}
	}
}
