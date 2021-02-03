package utils

import (
	"io/ioutil"
	"log"
	"os"
)

//PathExists determines whether the path to resource exists on disk
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//CreateFolder ...
func CreateFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil //path already exists
}

func writeFile(dst string, d []byte) error {
	//fmt.Printf("WriteFile: Size of download: %d\n", len(d))
	err := ioutil.WriteFile(dst, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
