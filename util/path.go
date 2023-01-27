package util

import (
	"io/ioutil"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateFileIfNotExist(filename string, content string) (bool, error) {
	pathExist, err := PathExists(filename)
	if err != nil {
		return false, err
	}
	if !pathExist {
		err := ioutil.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
