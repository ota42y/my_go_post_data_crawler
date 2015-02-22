package util

import (
	"io/ioutil"
)

func LoadFile(path string) (buf []byte) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return
}

