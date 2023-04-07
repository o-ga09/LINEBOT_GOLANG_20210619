package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"log"
	"main/store"
	"main/util"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	os.Setenv("PORT","8080")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set.")
	}

	listener, err := net.Listen("tcp",fmt.Sprintf(":%s",port))
	if err != nil {
		log.Fatal("port must be setten : ",err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthcheck)
	mux.HandleFunc("/callback", returnHTTPServer)

	server := &http.Server{
		Handler: mux,
	}

	go server.Serve(listener)

	quit := make(chan os.Signal,1)
	signal.Notify(quit,os.Interrupt)
	<- quit
	log.Printf("stopping server")
	server.Shutdown(context.Background())
}

func returnHTTPServer(w http.ResponseWriter, req *http.Request) {
	util.LoggingSettings("/usr/local/linebot01/test.log")
	client := &http.Client{
				Transport: &http.Transport{
				TLSHandshakeTimeout: 100 * time.Second,
				},
			}

	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"),linebot.WithHTTPClient(client))
	var reply_message string
	var hash string
	if err != nil {
		log.Fatalf("1")
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
		hash = util.Tohash(event.Source.UserID)
		isExit_master := store.Select_kaiin_host(hash)
		if len(isExit_master) == 0 {
			user, _ := bot.GetProfile(event.Source.UserID).Do()
			store.Insert_kaiin_host(hash, user.DisplayName)
			reply_message = "初めてですね！ありがとうございます。"
			if _, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage(reply_message)).Do(); err != nil {
				log.Print(err)
			}
		}
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.Contains(message.Text, "おはよう") || strings.Contains(message.Text, "こんにちは") || strings.Contains(message.Text, "こんばんは") {
					reply_message = greeting(message.Text)
				}else if message.Text == "家計簿" {
					reply_message = "https://liff.line.me/1656121916-RkvbrzoM"
				}else if message.Text == "健康管理" || strings.Contains(message.Text, "答") {
					reply_message = bodymanagement(hash, message.Text)
				} else {
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
					linebot.NewQuickReplyButton("", linebot.NewMessageAction("株価", "株価")),
					linebot.NewQuickReplyButton("", linebot.NewMessageAction("健康管理", "健康管理")),
					linebot.NewQuickReplyButton("", linebot.NewMessageAction("ぐるなび検索", "ぐるなび")),
					linebot.NewQuickReplyButton("", linebot.NewMessageAction("山本彩", "山本彩")),
				),
			)

			if _, err := bot.PushMessage(event.Source.UserID, resp).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type","application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
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

func bodymanagement(user_id string, text string) string {
	reply_message := ""
	result_bodyinfo := store.Select_bodymanagement(user_id)
	if text == "健康管理" {
		if len(result_bodyinfo) == 0 {
			reply_message = "初めてですね\n身長と体重を教えてください\n例)答 [身長] [体重]"
		} else {
			reply_message = "今日の体重を教えてね\n"
		}
	} else {
		tmp := strings.Split(text, " ")
		if len(tmp) == 3 {
			store.Insert_bodymanagement(user_id, tmp[1], tmp[2])
			reply_message = "正しく記録したよ！"
		} else if len(tmp) == 2 {
			height := strconv.FormatFloat(result_bodyinfo[0].Height, 'f', -1, 64)
			store.Insert_bodymanagement(user_id, tmp[1], height)
			reply_message = "正しく記録したよ！"
		}
	}
	return reply_message
}