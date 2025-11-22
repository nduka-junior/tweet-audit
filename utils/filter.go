package utils
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func FilterTweets( filePath string) (string , error) {
	// filePath := "tweets_clean.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return "", err
	}

	var tweets []map[string]interface{}
	err = json.Unmarshal(data, &tweets)
	if err != nil {
		fmt.Printf("JSON parse error: %v\n", err)
		return "", err
	}

	// store filtered tweets as slice of string pairs
	var filteredTweets []string

	for _, item := range tweets {
		tweet := item["tweet"].(map[string]interface{})

		fullText := tweet["full_text"].(string)

		// Use id_str if available; fallback to id
		var tweetID string
		if idStr, ok := tweet["id_str"].(string); ok {
			tweetID = idStr
		} else if idFloat, ok := tweet["id"].(float64); ok {
			tweetID = fmt.Sprintf("%.0f", idFloat)
		} else {
			tweetID = "unknown"
		}

		url := "https://x.com/i/web/status/" + tweetID

		// append text + url pair
		filteredTweets = append(filteredTweets, fullText, url)
	}

	// Prepare CSV rows
	var rows [][]string
	rows = append(rows, []string{"text", "url"}) // header

	for i := 0; i < len(filteredTweets); i += 2 {
		text := filteredTweets[i]
		url := filteredTweets[i+1]
		rows = append(rows, []string{text, url})
	}

	// Write CSV
	file, err := os.Create("tweets.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(rows)

	fmt.Println("CSV created successfully!")
	return "tweets.csv",nil
}
