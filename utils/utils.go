package utils

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"strings"
)

// RandToken returns a random token
// it must be remplaced by UUID
func RandToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// FindRealFilename tries to find a file associated to an id
// if the file exists, his entire name is returned
func FindRealFilename(path, id string) string {
	items, _ := ioutil.ReadDir(path)
	for _, item := range items {
		if strings.HasPrefix(item.Name(), id) {
			return path + item.Name()
		}
	}

	return ""
}
