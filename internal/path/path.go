package path

import (
	"os"
	"regexp"
)

const projectDirName = "write-my-naita"

func GetProjectRootPath() string {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	return string(rootPath)
}
