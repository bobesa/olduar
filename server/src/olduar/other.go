package olduar

import (
	"io/ioutil"
)

func GetFilesFromDirectory(path string) []string {
	filePaths := make([]string, 0)

	files, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range files {
			if file.IsDir() {
				filePaths = append(filePaths, GetFilesFromDirectory(path+"/"+file.Name())...)
			} else {
				filePaths = append(filePaths, path+"/"+file.Name())
			}
		}
	}

	return filePaths
}
