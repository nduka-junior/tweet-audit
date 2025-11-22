package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"github.com/nduka-junior/tweet-audit/utils"
)

// TweetResult handles both string and boolean flag output
type TweetResult struct {
	TweetText      string      `json:"tweet_text"`
	TweetURL       string      `json:"tweet_url"`
	FlagForDeletion interface{} `json:"flag_for_deletion"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not loaded")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Convert JS to JSON
	jsonPath,err := utils.ConvertJsFileToJSON("tweets.js", "tweets_clean.json")
	if err != nil {
		log.Fatal("Failed to convert JS to JSON:", err)
	}
	if jsonPath == "" {
		log.Fatal("Failed to convert JS to JSON")
	}
	// Filter tweets and create CSV
	csvPath, err := utils.FilterTweets(jsonPath)
	if err != nil {
		log.Fatal("Failed to filter tweets:", err)
	}
	// Open CSV with tweets
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Skip header row
	tweets := records[1:]

	batchSize := 50
	for i := 0; i < len(tweets); i += batchSize {
		end := i + batchSize
		if end > len(tweets) {
			end = len(tweets)
		}
		batch := tweets[i:end]

		// Build batch prompt
		var batchText []string
		for j, record := range batch {
			batchText = append(batchText, fmt.Sprintf("%d. Text: \"%s\" URL: %s", j+1, record[0], record[1]))
		}

		prompt := fmt.Sprintf(`
You are an AI assistant tasked with auditing tweets for deletion.

Rules:
1. Flag a tweet for deletion if it contains unprofessional language, offensive content, or specific keywords.
2. Flag tweets that express outdated opinions or views no longer supported.
3. Return strictly JSON array ONLY, without any markdown or code blocks, with objects: tweet_text, tweet_url, flag_for_deletion.

Evaluate these tweets:
%s
`, strings.Join(batchText, "\n"))

		// Call Gemini
		result, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(prompt),
			nil,
		)
		if err != nil {
			log.Println("Error generating content:", err)
			continue
		}

		// Clean AI output from backticks or extra spaces
		clean := strings.TrimSpace(result.Text())
		clean = strings.TrimPrefix(clean, "```json")
		clean = strings.TrimPrefix(clean, "```")
		clean = strings.TrimSuffix(clean, "```")

		// Parse JSON
		var parsed []TweetResult
		err = json.Unmarshal([]byte(clean), &parsed)
		if err != nil {
			log.Println("Error parsing JSON:", err, "\nRaw output:", clean)
			continue
		}

		// Save to CSV
		saveResultsToCSV(parsed)
	}
}

func saveResultsToCSV(results []TweetResult) {
	file, err := os.OpenFile("tweets_flagged.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header only if file is empty
	info, _ := file.Stat()
	if info.Size() == 0 {
		writer.Write([]string{"tweet_text", "tweet_url", "flag_for_deletion"})
	}

	for _, r := range results {
		flag := fmt.Sprintf("%v", r.FlagForDeletion) // handles bool or string
		writer.Write([]string{r.TweetText, r.TweetURL, flag})
	}
	writer.Flush()
}
