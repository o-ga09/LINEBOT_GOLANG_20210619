package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/syumai/workers/cloudflare"
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
	client := &http.Client{}
	channel_secret := cloudflare.Getenv("LINE_CHANNEL_SECRET")
	access_token := cloudflare.Getenv("LINE_ACCESS_TOKEN")
	bot, err := messaging_api.NewMessagingApiAPI(access_token, messaging_api.WithHTTPClient(client))

	var reply_message string
	if err != nil {
		log.Fatalf("can not connect line messaging api: %v", err)
	}

	cb, err := webhook.ParseRequest(channel_secret, req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	var name, content string
	for _, event := range cb.Events {
		// eventの種類でスイッチ
		switch e := event.(type) {
		case webhook.MessageEvent:
			// profileを取得
			switch s := e.Source.(type) {
			case webhook.UserSource:
				res, err := bot.GetProfile(s.UserId)
				if err != nil {
					log.Println(err)
				}
				name = res.DisplayName
			}
			// 送信メッセージを取得
			switch m := e.Message.(type) {
			case webhook.TextMessageContent:
				content = m.Text
			}

			// メッセージを返信する
			reply_message = fmt.Sprintf("%sさん！ありがとうございます。", name)
			if _, err := bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
				ReplyToken: e.ReplyToken,
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{
						Text: reply_message,
					},
					messaging_api.TextMessage{
						Text: content,
					},
				},
			}); err != nil {
				log.Print(err)
			}
		}
	}
}
