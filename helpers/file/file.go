package file

import (
	"os"
)

func FileExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}

func IsFileReadable(path string) (bool) {
	file, err := os.Open(path);

	if err != nil {
		return false
	} else {
		err = file.Close()
		if err != nil {
			return false
		}
	}

	return true
}
