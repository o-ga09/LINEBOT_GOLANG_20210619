package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"html/template"
	"io"
	"log"
	"main/myfunc"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type couponURL struct {
	Pc     string `json:"pc"`
	Mobile string `json:"mobile"`
}

type imageURL struct {
	ShopImage1 string `json:"shop_image1"`
	ShopImage2 string `json:"shop_image2"`
	QRcode     string `json:"qrcode"`
}

type access struct {
	Line        string      `json:"line"`
	Station     string      `json:"station"`
	StationExit string      `json:"station_exit"`
	Walk        json.Number `json:"walk"`
	Note        string      `json:"note"`
}

type pr struct {
	PRShort string `json:"pr_short"`
	PRLong  string `json:"pr_long"`
}

type code struct {
	Areacode      string   `json:"reacode"`
	Areaname      string   `json:"areaname"`
	Prefcode      string   `json:"prefcode"`
	Prefname      string   `json:"prefname"`
	AreacodeS     string   `json:"areacode_s"`
	AreanameS     string   `json:"areaname_s"`
	CategoryCodeL []string `json:"category_code_l"`
	CategoryNameL []string `json:"category_name_l"`
	CategoryCodeS []string `json:"category_code_s"`
	CategoryNameS []string `json:"category_name_s"`
}

type flag struct {
	MobileSite   json.Number `json:"mobile_site"`
	MobileCoupon json.Number `json:"mobile_coupon"`
	PcCoupon     json.Number `json:"pc_coupon"`
}

type rest struct {
	Attributes struct {
		Order json.Number `json:"order"`
	} `json:"@attributes"`
	ID          string      `json:"id"`
	Updatedate  time.Time   `json:"update_date"`
	Name        string      `json:"name"`
	NameKana    string      `json:"name_kana"`
	Latitude    json.Number `json:"latitude"`
	Longitude   json.Number `json:"longitude"`
	Category    string      `json:"category"`
	URL         string      `json:"url"`
	URLMobile   string      `json:"url_mobile"`
	CouponURL   couponURL   `json:"coupon_url"`
	ImageURL    imageURL    `json:"image_url"`
	Address     string      `json:"address"`
	Tel         string      `json:"tel"`
	TelSub      string      `json:"tel_sub"`
	Fax         string      `json:"fax"`
	OpenTime    string      `json:"opentime"`
	Holiday     string      `json:"holiday"`
	Access      access      `json:"access"`
	ParkingLots json.Number `json:"parking_lots"`
	PR          pr          `json:"pr"`
	Code        code        `json:"code"`
	Budget      json.Number `json:"budget"`
	Party       json.Number `json:"party"`
	Lunch       json.Number `json:"lunch"`
	CreditCard  string      `json:"credit_card"`
	EMoney      string      `json:"e_money"`
	Flag        flag        `json:"flags"`
}

type getJson struct {
	Attributes struct {
		APIversion string `json:"api_version"`
	} `json:"@attributes"`
	Total_hit_count json.Number `json:"total_hit_count"`
	Hit_per_page    json.Number `json:"hit_per_page"`
	Page_offset     json.Number `json:"page_offset"`
	Rest            []rest      `json:"rest"`
}

func LoggingSettings(logFile string) {
    //_=error
    //os.O_RDWR　READ　WRITE 読み書き両方する時
    //os.O_CREATE 存在しなかった場合新規ファイルを作成する場合
    //os.O_APPEND  ファイルに追記したいとき
    //0666 
    //引数: ファイルのパス, フラグ, パーミッション(わからなければ0666でおっけーです)
    //上記モード指定。読み込む、作成、権限(０６６６＝読み書き）を設定。
    logfile, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    //stdout 画面上に出る出力 をlogfileに書き込む
    multiLogFile := io.MultiWriter(os.Stdout, logfile)
    //フォーマット指定
    //日付、時間、短いエラーの名前
    log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
    //ログファイルの出力先を変更   
    log.SetOutput(multiLogFile)
}

func redirect(w http.ResponseWriter, req *http.Request) {

    target := "https://" + req.Host + req.URL.Path 
    if len(req.URL.RawQuery) > 0 {
        target += "?" + req.URL.RawQuery
    }
    log.Printf("redirect to: %s", target)
    http.Redirect(w, req, target,301)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set.")
	}

	go http.ListenAndServe(":80",http.HandlerFunc(redirect))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	mux.HandleFunc("/", Handler)
	mux.HandleFunc("/cashform", cashformHandler)
	mux.HandleFunc("/weightform_1", weightform_firstHandler)
	mux.HandleFunc("/weightform_2", weightform_secondHandler)
	//mux.HandleFunc("/", HelloHandler)
	mux.HandleFunc("/callback", returnHTTPServer)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               tls.VersionTLS12,
		},
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServeTLS("/etc/letsencrypt/live/ilovegameandcomputer.tk/fullchain.pem","/etc/letsencrypt/live/ilovegameandcomputer.tk/privkey.pem"); err != nil {
		log.Fatal(err)
	}
}
func returnHTTPServer(w http.ResponseWriter, req *http.Request) {
	LoggingSettings("/usr/local/linebot01/test.log")
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
		hash = Tohash(event.Source.UserID)
		isExit_master := myfunc.Select_kaiin_host(hash)
		if len(isExit_master) == 0 {
			user, _ := bot.GetProfile(event.Source.UserID).Do()
			myfunc.Insert_kaiin_host(hash, user.DisplayName)
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

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Helloworld!@@@@@@")
}

func Handler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("template/index.html"))
    tmpl.Execute(w, nil)
} 

func cashformHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("template/input_cash.html"))
    tmpl.Execute(w, nil)
} 

func weightform_firstHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("template/input_weight_atfirst.html"))
    tmpl.Execute(w, nil)
} 

func weightform_secondHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("template/input_weight_second.html"))
    tmpl.Execute(w, nil)
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
	result_bodyinfo := myfunc.Select_bodymanagement(user_id)
	if text == "健康管理" {
		if len(result_bodyinfo) == 0 {
			reply_message = "初めてですね\n身長と体重を教えてください\n例)答 [身長] [体重]"
		} else {
			reply_message = "今日の体重を教えてね\n"
		}
	} else {
		tmp := strings.Split(text, " ")
		if len(tmp) == 3 {
			myfunc.Insert_bodymanagement(user_id, tmp[1], tmp[2])
			reply_message = "正しく記録したよ！"
		} else if len(tmp) == 2 {
			height := strconv.FormatFloat(result_bodyinfo[0].Height, 'f', -1, 64)
			myfunc.Insert_bodymanagement(user_id, tmp[1], height)
			reply_message = "正しく記録したよ！"
		}
	}
	return reply_message
}

func getAPI(name string) getJson {
	grunavi_api_key := os.Getenv("GRUNAVI_API_KEY")
	url := "https://api.gnavi.co.jp/RestSearchAPI/v3/?keyid=" + grunavi_api_key + "&name=" + name
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data getJson

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	fmt.Println(data.Hit_per_page.Float64())
	// for _, item := range data {
	// 	fmt.Printf("%s %s\n", item.Total_hit_count, item.Hit_per_page)
	// }

	return data
}

func Tohash(originalstr string) string {
	salt := os.Getenv("SALT")
	hashstr := []byte(originalstr + salt)

	hash_sha256 := sha256.Sum256(hashstr)

	return hex.EncodeToString(hash_sha256[:])
}

func isMatch(hash,userid string) bool{
	salt := os.Getenv("SALT")
	hashstr := []byte(userid + salt)
	hash_sha256 := sha256.Sum256(hashstr)

	if hash == hex.EncodeToString(hash_sha256[:]) {
		return true
	}
	return false
}
