package main

import (
	"bufio"
	"fmt"
	"lakpahana/write-my-naita/internal/docx"
	"lakpahana/write-my-naita/internal/llm/gemini"
	"log"
	"os"
	"regexp"
	"strings"

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

	prompt := collectDailyProgress()

	report, err := GeminiClient.GenerateTimelineForAWeek(prompt)

	if err != nil {
		log.Fatalf("Error generating weekly training report: %v", err)
	}

	err = docx.InsertWeeklyTimelineToDocx(report, "output.docx")

	if err != nil {
		log.Fatalf("Error writing to docx: %v", err)
	}

	log.Printf("report: %v", report.DailyProgress.Day1)
}

type WeeklyProgress struct {
	Project   string
	Days      map[int]string
	Learnings string
}

func collectDailyProgress() string {
	reader := bufio.NewReader(os.Stdin)
	progress := WeeklyProgress{
		Days: make(map[int]string),
	}

	// Clear screen and show welcome message
	fmt.Print("\033[H\033[2J")
	fmt.Println("📝 Weekly Progress Tracker")
	fmt.Println("-------------------------")

	// Collect project details
	fmt.Println("\n🎯 First, let's capture your project details:")
	fmt.Print("\nProject Details(Name,Tech): ")
	progress.Project, _ = reader.ReadString('\n')

	// Collect daily progress
	fmt.Println("\n📅 Now, let's record your daily progress:")
	for day := 1; day <= 5; day++ {
		fmt.Printf("\nDay %d Progress (What did you accomplish?): ", day)
		progress.Days[day], _ = reader.ReadString('\n')
	}

	// Collect key learnings
	fmt.Println("\n💡 Finally, what are your key learnings from this week?")
	progress.Learnings, _ = reader.ReadString('\n')

	// Generate and return the prompt
	return generatePrompt(progress)
}

func generatePrompt(progress WeeklyProgress) string {
	return fmt.Sprintf(`
		Project Details: %s

		Daily Progress:
		Day 1: %s
		Day 2: %s
		Day 3: %s
		Day 4: %s
		Day 5: %s

		Key Learnings: %s
		`,
		strings.TrimSpace(progress.Project),
		strings.TrimSpace(progress.Days[1]),
		strings.TrimSpace(progress.Days[2]),
		strings.TrimSpace(progress.Days[3]),
		strings.TrimSpace(progress.Days[4]),
		strings.TrimSpace(progress.Days[5]),
		strings.TrimSpace(progress.Learnings))
}
