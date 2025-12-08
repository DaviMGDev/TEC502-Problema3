package utils 

import (
	"encoding/json"
)

type Dict map[string]any 

func (dict Dict) Json() ([]byte, error) {
	return json.Marshal(dict)
}

func (dict Dict) String() (string, error) {
	data, err := dict.Json()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

