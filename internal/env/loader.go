package env

import (
	"log"
	"os"

	"lakpahana/write-my-naita/internal/path"

	"github.com/joho/godotenv"
)

const projectDirName = "write-my-naita"

// LoadEnv loads environment variables from a `.env` file
func LoadEnv() {
	rootPath := path.GetProjectRootPath()
	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(-1)
	}
}
