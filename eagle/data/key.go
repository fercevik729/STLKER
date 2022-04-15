package data

import (
	"io/ioutil"
	"os"
)

func LoadKey(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	key, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(key), nil
}
