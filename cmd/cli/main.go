package main

import (
	"lakpahana/write-my-naita/internal/llm/gemini"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "write-my-naita"

// LoadEnv loads env vars from .env
func LoadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(-1)
	}
}

func main() {
	LoadEnv()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalf("API_KEY environment variable not set")
	}

	GeminiClient, err := gemini.NewGeminiClient(apiKey)
	if err != nil {
		log.Fatalf("Error creating Gemini client: %v", err)
	}

	prompt := `my project is to create a web application that can be used to manage a small business. The application should be able to store information about customers, products, and orders. It should also allow users to view reports and generate invoices.
	in this week i worked on authentication
	`
	report, err := GeminiClient.GenerateTimelineForAWeek(prompt)

	if err != nil {
		log.Fatalf("Error generating weekly training report: %v", err)
	}

	log.Printf("report: %v", report.DailyProgress.Day1)
}
