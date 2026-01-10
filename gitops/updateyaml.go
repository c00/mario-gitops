package gitops

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spyzhov/ajson"
	"go.yaml.in/yaml/v3"
)

// UpdateYaml will update a yaml file according a to jsonpath query and a new value.
// It will turn yaml into json into ajson, perform jsonpath query, update all found nodes, turn back into json into yaml
func updateYaml(input, path string, newVal any) (output string, err error) {
	var data interface{}

	err = yaml.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", err
	}

	jsonInput, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("cannot marshal json: %w", err)
	}

	//do the magic
	root, _ := ajson.Unmarshal(jsonInput)
	nodes, _ := root.JSONPath(path)

	for _, node := range nodes {
		node.Set(newVal)
	}
	result, _ := ajson.Marshal(root)

	err = json.Unmarshal(result, &data)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshall json back: %w", err)
	}

	// Marshal back to YAML
	b := bytes.Buffer{}
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	err = enc.Encode(data)
	if err != nil {
		return "", fmt.Errorf("cold not marshall data: %w", err)
	}

	return b.String(), nil
}
