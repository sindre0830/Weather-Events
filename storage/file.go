package storage

import (
	"io/ioutil"
	"os"
)

// FindCollection checks if file exist.
func FindCollection(filename string) bool {
	_, err := os.Stat("./data/" + filename + ".json")
	return err == nil
}

// ReadCollection reads file.
func ReadCollection(filename string) ([]byte, error) {
	return ioutil.ReadFile("./data/" + filename + ".json")
}

// WriteCollection writes to file.
func WriteCollection(filename string, file []byte) error {
	return ioutil.WriteFile("./data/"+filename+".json", file, 0644)
}
