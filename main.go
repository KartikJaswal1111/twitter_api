package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

const (
	tweetPostURL   = "https://api.twitter.com/1.1/statuses/update.json"
	tweetDeleteURL = "https://api.twitter.com/1.1/statuses/destroy/%s.json"
)

func main() {

	errone := godotenv.Load()
	if errone != nil {
		log.Fatal("Error loading .env file")
	}

	APIkey := os.Getenv("API_KEY")
	APIkeySecret := os.Getenv("API_Key_Secret")
	accessToken := os.Getenv("Access_Token")
	accessTokenSecret := os.Getenv("Access_Token_Secret")
	print(APIkey, APIkeySecret, accessToken, accessTokenSecret)

	config := oauth1.NewConfig(APIkey, APIkeySecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	tweetText := "This is my first tweet as the twitter user"
	err := postTweet(httpClient, tweetText)
	if err != nil {
		fmt.Println("Error posting tweet:", err)
		return
	}
	fmt.Println("Tweet posted successfully!")

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
