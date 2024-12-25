package llm

type WeeklyTrainingReport struct {
	DailyProgress struct {
		Day1 []string `json:"day_1"`
		Day2 []string `json:"day_2"`
		Day3 []string `json:"day_3"`
		Day4 []string `json:"day_4"`
		Day5 []string `json:"day_5"`
	} `json:"daily_progress"`
	KeyLearnings string `json:"key_learnings"`
}

type LLMProvider interface {
	GenerateTimelineForAWeek(prompt string) (WeeklyTrainingReport, error)
	// GenerateTimeline(prompt string, durationInWeeks int) (string, error)
}
