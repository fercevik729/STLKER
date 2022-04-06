package data

import (
	"encoding/json"
	"io"
)

func ToJSON(i interface{}, w io.Writer) {
	e := json.NewEncoder(w)
	e.Encode(i)
}

func FromJSON(i interface{}, r io.Reader) {
	d := json.NewDecoder(r)
	d.Decode(i)
}
