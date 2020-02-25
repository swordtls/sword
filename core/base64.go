package core

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func Base64P12(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	encoded := base64.StdEncoding.EncodeToString(file)
	fmt.Println(encoded)
}
