package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dghubble/oauth1"
)

const (
	tweetPostURL   = "https://api.twitter.com/1.1/statuses/update.json"
	tweetDeleteURL = "https://api.twitter.com/1.1/statuses/destroy/%s.json"
)

func main() {
	// Parse command-line arguments for authorization
	if len(os.Args) != 5 {
		fmt.Println("Usage: ./your_program <consumer_key> <consumer_secret> <access_token> <access_token_secret>")
		os.Exit(1)
	}

	consumerKey := os.Args[1]
	consumerSecret := os.Args[2]
	accessToken := os.Args[3]
	accessTokenSecret := os.Args[4]

	// Create OAuth config
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Post a new tweet
	tweetText := "This is my first tweet as the twitter user" // Replace with your desired tweet
	err := postTweet(httpClient, tweetText)
	if err != nil {
		fmt.Println("Error posting tweet:", err)
		return
	}
	fmt.Println("Tweet posted successfully!")

	// Delete a tweet (replace with a valid tweet ID)
	tweetID := "<YOUR_TWEET_ID>" // Enter the ID of the tweet you want to delete
	err = deleteTweet(httpClient, tweetID)
	if err != nil {
		fmt.Println("Error deleting tweet:", err)
		return
	}
	fmt.Println("Tweet deleted successfully!")
}

func postTweet(client *http.Client, text string) error {
	req, err := http.NewRequest("POST", tweetPostURL,
		strings.NewReader(fmt.Sprintf("status=%s", text)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error posting tweet: %s", string(body))
	}

	var tweetData map[string]interface{}
	err = json.Unmarshal(body, &tweetData)
	if err != nil {
		return fmt.Errorf("Error parsing tweet response: %s", err)
	}

	fmt.Println("Tweet created: https://twitter.com/status/" + tweetData["id_str"].(string))
	return nil
}

func deleteTweet(client *http.Client, tweetID string) error {
	url := fmt.Sprintf(tweetDeleteURL, tweetID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error deleting tweet: %s", err)
		}
		return fmt.Errorf("Error deleting tweet: %s", string(body))
	}

	fmt.Println("Tweet deleted successfully!")
	return nil
}
