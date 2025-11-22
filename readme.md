Here’s the updated README with your requested reference included:

---

# Tweet-Audit: AI-Powered Twitter/X Archive Analysis

**Tweet-Audit** is a Go-based tool that helps you analyze your Twitter (X) archive and flag tweets for potential deletion based on custom criteria. Using Google’s Gemini AI, the tool evaluates your tweets for unprofessional language, offensive content, outdated opinions, or other user-defined rules, and outputs flagged tweets in a CSV format.

> **Note:** This project was initially drafted by **benx421** to help improve backend knowledge and serve as a learning tool for building Go-based data pipelines and AI integrations.

---

## **Features**

* Convert exported Twitter/X JS archive to JSON.
* Extract tweets and generate a CSV of tweet text and URLs.
* Batch auditing using Google Gemini AI.
* Flag tweets automatically based on defined rules.
* Save audit results into `tweets_flagged.csv` for easy review.

---

## **Requirements**

* Go 1.21+
* `.env` file with your Gemini API key:

```
GEMINI_API_KEY=your_api_key_here
```

* Twitter/X archive (JS file exported from X).

---

## **Installation**

1. Clone the repository:

```bash
git clone https://github.com/nduka-junior/tweet-audit.git
cd tweet-audit
```

2. Initialize Go modules:

```bash
go mod tidy
```

3. Create a `.env` file with your Gemini API key:

```bash
touch .env
echo "GEMINI_API_KEY=your_api_key_here" >> .env
```

---

## **Project Structure**

```
tweet-audit/
├── go.mod
├── main.go                  # Main execution file
├── utils/
│   ├── jsToJSON.go          # Converts Twitter JS archive to JSON
│   └── filter.go            # Filters tweets and generates CSV
├── tweets.js                # Your exported Twitter/X JS archive
├── tweets_clean.json        # Generated JSON file
├── tweets.csv               # Extracted tweets CSV
└── tweets_flagged.csv       # Flagged tweets output
```

---

## **Usage**

1. Place your Twitter/X archive `tweets.js` in the project folder.
2. Run the program:

```bash
go run main.go
```

3. The tool performs:

   * JS → JSON conversion (`tweets_clean.json`)
   * Extract tweets into CSV (`tweets.csv`)
   * Batch audit tweets with Gemini AI
   * Save flagged tweets in `tweets_flagged.csv`

4. **Example flagged CSV output:**

| tweet_text         | tweet_url                                                            | flag_for_deletion |
| ------------------ | -------------------------------------------------------------------- | ----------------- |
| Example tweet text | [https://x.com/i/web/status/12345](https://x.com/i/web/status/12345) | true              |
| Another tweet      | [https://x.com/i/web/status/67890](https://x.com/i/web/status/67890) | false             |

---

## **Custom Rules**

You can modify the AI prompt in `main.go` to adjust the audit rules. For example:

* Flag unprofessional language
* Flag offensive content
* Flag outdated opinions
* Include any specific keywords you want to monitor

---

## **Batch Processing**

* The tool processes tweets in batches (default: 50) to reduce API cost and improve performance.
* Gemini AI evaluates each batch and returns JSON results.

---

## **Dependencies**

* [Go](https://golang.org/)
* [joho/godotenv](https://github.com/joho/godotenv) — for environment variables
* [Google GenAI Go SDK](https://pkg.go.dev/google.golang.org/genai) — for AI text generation

Install dependencies:

```bash
go get github.com/joho/godotenv
go get google.golang.org/genai
```

---

## **Notes**

* Make sure your Gemini API key is valid.
* Ensure the Twitter JS archive is downloaded from your account settings.
* The `tweets_flagged.csv` is appended on every run; remove it if you want a fresh audit.

---

## **License**

MIT License © 2025 Nduka Sochi John-Junior
