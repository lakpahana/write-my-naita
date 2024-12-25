package cli

import (
	"fmt"
	"lakpahana/write-my-naita/internal/env"
	"lakpahana/write-my-naita/internal/llm/gemini"
	"os"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a training report",
	Long:  `Generate a weekly training report using an LLM.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load environment variables
		env.LoadEnv()

		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			fmt.Println("Error: API_KEY environment variable not set")
			return
		}

		geminiClient, err := gemini.NewGeminiClient(apiKey)

		if err != nil {
			fmt.Println("Error creating Gemini client")
			return
		}

		prompt := "Write a weekly training report for the prompt."

		result, err := geminiClient.GenerateTimelineForAWeek(prompt)

		if err != nil {
			fmt.Println("Error generating report")
			return
		}

		fmt.Println("Generated Report:")
		fmt.Println(result)
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
}
