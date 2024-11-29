package utils

import "encoding/json"

func DumpJSON(v any) ([]byte, error) {

	jsonBytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return []byte{}, err
	}

	return jsonBytes, nil
}
