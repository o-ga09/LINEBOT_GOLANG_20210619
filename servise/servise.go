package servise

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func ResponseBot(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{
				Transport: &http.Transport{
				TLSHandshakeTimeout: 100 * time.Second,
				},
			}

	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"),linebot.WithHTTPClient(client))
	var reply_message string
	if err != nil {
		log.Fatalf("can not connect line messaging api")
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		res := bot.GetProfile(event.Source.UserID)
		profile, err := res.Do()
		name := profile.DisplayName
		if err != nil { 
			log.Printf("can not get profile")
			name = "名無しさん"
		}
		reply_message = fmt.Sprintf("%sさん！ありがとうございます。",name)
		if _, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage(reply_message)).Do(); err != nil {
			log.Print(err)
		}
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.Contains(message.Text, "おはよう") || strings.Contains(message.Text, "こんにちは") || strings.Contains(message.Text, "こんばんは") {
					reply_message = greeting(message.Text)
				}else {
					reply_message = message.Text
				}

				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply_message)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf("Sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
			resp := linebot.NewTextMessage(
				"用件は何？",
			).WithQuickReplies(
				linebot.NewQuickReplyItems(
					linebot.NewQuickReplyButton("", linebot.NewMessageAction("天気", "天気")),
					linebot.NewQuickReplyButton("", linebot.NewMessageAction("山本彩", "山本彩")),
				),
			)

			if _, err := bot.PushMessage(event.Source.UserID, resp).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}

func greeting(message string) string {
	var timezoneJST = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.Local = timezoneJST
	time.LoadLocation("Asia/Tokyo")

	now := time.Now()
	if now.Hour() >= 5 && now.Hour() <= 10 {
		return "おはよう"
	} else if now.Hour() >= 11 && now.Hour() <= 16 {
		return "こんにちは"
	} else if now.Hour() >= 17 && now.Hour() <= 23 {
		return "こんばんは"
	}
	return "おやすみ"
}