package utils

import "encoding/json"

type Dict map[string]any

func (dict Dict) Json() ([]byte, error) {
	return json.Marshal(dict)
}

func (dict Dict) String() (string, error) {
	bytes, err := dict.Json()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJson(data []byte) (Dict, error) {
	var dict Dict
	err := json.Unmarshal(data, &dict)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

func FromString(data string) (Dict, error) {
	return FromJson([]byte(data))
}
