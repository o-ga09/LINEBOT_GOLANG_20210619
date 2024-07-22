package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func (h *Handler) CallBack(w http.ResponseWriter, req *http.Request) {
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))

	var reply_message string
	if err != nil {
		log.Fatalf("can not connect line messaging api: %v", err)
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
			log.Fatal(err)
		}
		reply_message = fmt.Sprintf("%sさん！ありがとうございます。", name)
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply_message)).Do(); err != nil {
			log.Print(err)
		}
	}
}
