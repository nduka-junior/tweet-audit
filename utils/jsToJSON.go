package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ConvertJsFileToJSON(inputPath string, outputPath string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return "", err
	}

	content := string(data)
	index := strings.Index(content, "=")
	if index == -1 {
		return "", fmt.Errorf("'=' not found in file")
	}

	jsonPart := strings.TrimSpace(content[index+1:])
	jsonPart = strings.TrimSuffix(jsonPart, ";")

	var tweets []map[string]interface{}
	err = json.Unmarshal([]byte(jsonPart), &tweets)
	if err != nil {
		return "", fmt.Errorf("JSON parse error: %v", err)
	}

	pretty, err := json.MarshalIndent(tweets, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %v", err)
	}

	err = os.WriteFile(outputPath, pretty, 0644)
	if err != nil {
		return "", fmt.Errorf("File write error: %v", err)
	}

	fmt.Println("Success! Saved JSON to:", outputPath)
	return outputPath, nil
}
