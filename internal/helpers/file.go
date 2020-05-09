package helpers

import "io/ioutil"

//LoadFile load file data
func LoadFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
