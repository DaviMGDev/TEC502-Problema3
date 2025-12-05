package utils

import (
	"encoding/json"
)

type Dict map[any]any 

func (dict Dict) String() string {
	data, _ := json.MarshalIndent(dict, "", "  ")
	return string(data)
}

func (dict Dict) Json() []byte {
	data, _ := json.MarshalIndent(dict, "", "  ")
	return data
}
