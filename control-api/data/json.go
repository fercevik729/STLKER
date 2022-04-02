package data

import (
	"encoding/json"
	"io"
)

func ToJSON(i interface{}, w io.Writer) {
	e := json.NewEncoder(w)
	e.Encode(i)
}
