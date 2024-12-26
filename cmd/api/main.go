package main

import (
	"encoding/json"
	"fmt"
	"lakpahana/write-my-naita/internal/docx"
	"lakpahana/write-my-naita/internal/env"
	"lakpahana/write-my-naita/internal/llm/gemini"
	"lakpahana/write-my-naita/internal/path"
	"log"
	"net/http"
	"os"
	"time"
)

type WeeklyProgress struct {
	Project   string         `json:"project"`
	Days      map[int]string `json:"days"`
	Learnings string         `json:"learnings"`
}

func main() {

	os.MkdirAll("downloads", os.ModePerm)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.GetProjectRootPath()+"/frontend/index.html")
	})

	http.HandleFunc("/api/generate", handleGenerate)
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir("downloads"))))

	http.HandleFunc("/api/download", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		filepath := fmt.Sprintf("downloads/%s", filename)
		http.ServeFile(w, r, filepath)
	})

	handler := corsMiddleware(http.DefaultServeMux)

	fmt.Println("Server starting on :8080...\n Open http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var progress WeeklyProgress
	if err := json.NewDecoder(r.Body).Decode(&progress); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prompt := generatePrompt(progress)

	env.LoadEnv()

	apiKey := os.Getenv("API_KEY")
	geminiClient, err := gemini.NewGeminiClient(apiKey)
	if err != nil {
		http.Error(w, "Failed to initialize Gemini client", http.StatusInternalServerError)
		return
	}

	report, err := geminiClient.GenerateTimelineForAWeek(prompt)
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("report_%s.docx", time.Now().Format("20060102150405"))
	filepath := fmt.Sprintf("downloads/%s", filename)

	err = docx.InsertWeeklyTimelineToDocx(report, filepath)
	if err != nil {
		http.Error(w, "Failed to create DOCX", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"fileUrl": "api/download?filename=" + filename,
	})
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
		progress.Project,
		progress.Days[1],
		progress.Days[2],
		progress.Days[3],
		progress.Days[4],
		progress.Days[5],
		progress.Learnings)
}
