package data

import (
	"encoding/json"
	"io"
)

func ToJSON(i interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(i)
}
