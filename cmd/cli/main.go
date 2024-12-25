package main

import (
	"lakpahana/write-my-naita/internal/llm/gemini"
	"log"
)

func main() {
	GeminiClient, err := gemini.NewGeminiClient("SDAD")
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

	log.Printf("report: %v", report)

}
