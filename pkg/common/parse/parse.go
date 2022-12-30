package parse

import (
	"encoding/json"
	"io"
)

func JsonFromReader(source io.Reader, target any) error {
	bytes, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func JsonFromBytes(source []byte, target any) error {
	return json.Unmarshal(source, target)
}
