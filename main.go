package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"

	youtube "google.golang.org/api/youtube/v3"

	"github.com/skratchdot/open-golang/open"
)

func saveToken(token *oauth2.Token) {
	f, err := os.Create("token.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	jsonStr, _ := json.Marshal(token)
	fmt.Fprintln(f, string(jsonStr))
}

func getToken(conf *oauth2.Config) *oauth2.Token {
	f, err := os.Open("token.json")
	var r io.ReadCloser
	if err == nil {
		defer f.Close()
		var t oauth2.Token
		dec := json.NewDecoder(f)
		err = dec.Decode(&t)
		if err != nil {
			log.Fatal(err)
		}
		return &t
	}

	url := conf.AuthCodeURL("test")
	open.Run(url)
	r = os.Stdin

	fmt.Println("Please enter the token which is displayed in your browser after you accepted pollydent. You can enter the token by copy & paste.")

	var code string
	var sc = bufio.NewScanner(r)
	if sc.Scan() {
		code = sc.Text()
	}

	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	saveToken(tok)
	return tok
}

func getChatIDFromVideo(service *youtube.Service) string {
	vlc := service.Videos.List("snippet,contentDetails,liveStreamingDetails")
	fmt.Println("Please set video id by copy & paste:")
	var sc = bufio.NewScanner(os.Stdin)
	var videoID string
	if sc.Scan() {
		videoID = sc.Text()
	}
	vlc.Id(videoID)
	resp, err := vlc.Do()
	if err != nil {
		log.Fatal(err)
	}

	var chatID string
	for _, i := range resp.Items {
		chatID = i.LiveStreamingDetails.ActiveLiveChatId
		break
	}
	return chatID
}

func listChatMessages(service *youtube.Service, chatID string, pageToken string) string {
	lclc := service.LiveChatMessages.List(chatID, "snippet,authorDetails")
	if pageToken != "" {
		lclc.PageToken(pageToken)
	}
	msgs, err := lclc.Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range msgs.Items {
		jsonBytes, err := json.Marshal(i)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonBytes))
	}
	time.Sleep(time.Duration(msgs.PollingIntervalMillis) * time.Millisecond)

	return msgs.NextPageToken
}

func main() {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	tok := getToken(conf)
	if tok == nil {
		log.Fatal("getting a token was failed")
	}

	client := conf.Client(oauth2.NoContext, tok)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatal(err)
	}

	chatID := getChatIDFromVideo(service)
	if chatID == "" {
		log.Fatal("chat id could not be found in video. the specified video id is invalid, or chat does not exist in the video.")
	}

	var pageToken string
	for {
		pageToken = listChatMessages(service, chatID, pageToken)
	}
}
