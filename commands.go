package main

import (
	"log"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/codegangsta/cli"
	"github.com/darkhelmet/twitterstream"
	"github.com/ikawaha/kagome/dic"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/joho/godotenv"
)

var YURIKO_ID int64 = 2725436402

var Commands = []cli.Command{
	commandLoves,
}

var commandLoves = cli.Command{
	Name:  "loves",
	Usage: "",
	Description: `
`,
	Action: loves,
}

func loves(c *cli.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := twitterstream.NewClient(
		os.Getenv("TW_CONSUMER_KEY"),
		os.Getenv("TW_CONSUMER_SECRET"),
		os.Getenv("TW_CLIENT_KEY"),
		os.Getenv("TW_CLIENT_SECRET"),
	)
	for {
		conn, err := client.Follow(string(YURIKO_ID))
		if err != nil {
			time.Sleep(1 * time.Minute)
			continue
		}
		kanojoru(conn)
	}
}

func kanojoru(conn *twitterstream.Connection) {
	for {
		tweet, err := conn.Next()
		if err != nil {
			return
		}
		if !canKanojonize(tweet) {
			continue
		}

		text := tweet.Text
		kanojonized := kanojonize(text)
		post(kanojonized)
	}
}

func kanojonize(text string) string {
	t := tokenizer.NewTokenizer()
	morphs, err := t.Tokenize(text)
	if err != nil {
		return ""
	}

	var strings []string
	for _, val := range morphs {
		var curr dic.Content
		curr, err = val.Content()
		if err != nil {
			continue
		} else if (curr.Pos == "名詞") && (curr.Pos1 == "一般") {
			// 置き換える単語を保持
			strings = append(strings[:], val.Surface)
		}
	}
	kanojonized := replaceToHisaichi(strings, text)
	return kanojonized
}

func replaceToHisaichi(strings []string, text string) string {
	for _, val := range strings {
		text = regexp.MustCompile(val).ReplaceAllString(text, "ひさいち")
	}
	return text
}

func canKanojonize(tweet *twitterstream.Tweet) bool {
	if len(tweet.Entities.Mentions) >= 1 {
		// メンションがあるときは無視
		return false
	} else if tweet.User.Id != YURIKO_ID {
		// 吉高由里子以外無視
		return false
	}
	return true
}

func post(kanojonized string) {
	anaconda.SetConsumerKey(os.Getenv("TW_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TW_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(
		os.Getenv("TW_CLIENT_KEY"),
		os.Getenv("TW_CLIENT_SECRET"),
	)

	v := url.Values{}
	api.PostTweet(kanojonized, v)
}
