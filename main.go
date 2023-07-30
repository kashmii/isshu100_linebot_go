package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Waka struct {
	No      int    `json:"no"`
	Kami    string `json:"kami"`
	Simo    string `json:"simo"`
	Kami_kana    string `json:"kami_kana"`
	Simo_kana    string `json:"simo_kana"`
	Author    string `json:"sakusya"`
	Author_kana    string `json:"sakusya_name"`
	Image    string `json:"image"`
}

type Choice struct {
	Kami    string `json:"kami"`
	Simo    string `json:"simo"`
	Correct string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	// ハンドラの登録
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/callback", lineHandler)

	fmt.Println("http://localhost:8080 で起動中...")
	// HTTPサーバを起動
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!!!!"
	fmt.Fprintf(w, msg)
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	// BOTを初期化
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// ===============
	// JSONファイルを読み込む
	file, err := os.Open("whole_master.json")
	if err != nil {
		fmt.Println("JSONファイルをオープンできませんでした:", err)
		return
	}
	defer file.Close()

	// ファイル内容をバイト配列に読み込む
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("JSONファイルの読み込みに失敗しました:", err)
		return
	}

	// バイト配列を構造体に変換
	var masterData []Waka
	err = json.Unmarshal(data, &masterData)
	if err != nil {
		fmt.Println("JSONデータの解析に失敗しました:", err)
		return
	}
	t := reflect.TypeOf(masterData)
	fmt.Println("データ型:", t)

	// 和歌を4つ取得して配列：長さ3に入れる
	// 配列[0]は正解として、それ以外ははずれ選択肢として使う
	randomElements := GetRandomWakas(masterData, 4)
	choices := make([]Choice, 4)
	for i, ele := range randomElements {
		choices[i].Kami = ele.Kami
		choices[i].Simo = ele.Simo
		if i == 0 {
			choices[i].Correct = "正解"
		} else {
			choices[i].Correct = "はずれ"
		}
	}

	question := choices[0].Kami

	// 選択肢4つをランダムに並べ直し
	rand.Seed(time.Now().UnixNano()) //乱数のシード設定
	rand.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})

	// ===============
	// ===============



	// リクエストからBOTのイベントを取得
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// メッセージがテキスト形式の場合
			case *linebot.TextMessage:
				switch message.Text {
				case "flex":
					resp := linebot.NewFlexMessage(
						"this is a flex message",
						&linebot.BubbleContainer{
							Type: linebot.FlexContainerTypeBubble,
							Body: &linebot.BoxComponent{
								Type:   linebot.FlexComponentTypeBox,
								Layout: linebot.FlexBoxLayoutTypeVertical,
								Contents: []linebot.FlexComponent{
									&linebot.TextComponent{
										Type: linebot.FlexComponentTypeText,
										Text: "'" + question + "'",
										Weight: "regular",
										Size: "md",
										Align: "center",
									},
									&linebot.TextComponent{
										Type: "text",
										Text: "下の句はどれでしょう？",
										Weight: "regular",
										Size: "md",
										Align: "center",
										Margin: "md",
									},
								},
							},
							Footer: &linebot.BoxComponent{
								Type:   linebot.FlexComponentTypeBox,
								Layout: linebot.FlexBoxLayoutTypeVertical,
								Spacing: linebot.FlexComponentSpacingTypeSm,
								Contents: []linebot.FlexComponent{
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										Action: linebot.NewMessageAction(choices[0].Simo, choices[0].Correct),
									},
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										Action: linebot.NewMessageAction(choices[1].Simo, choices[1].Correct),
									},
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										Action: linebot.NewMessageAction(choices[2].Simo, choices[2].Correct),
									},
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										Action: linebot.NewMessageAction(choices[3].Simo, choices[3].Correct),
									},
								},
							},
						},
					)
					_, err = bot.ReplyMessage(event.ReplyToken, resp).Do()
					if err != nil {
						log.Print(err)
					}
				


				case "as":
				resp := linebot.NewFlexMessage(
					"this is a flex message",
					&linebot.BubbleContainer{
						Type: linebot.FlexContainerTypeBubble,
						Hero: &linebot.ImageComponent{
							Type:        "image",
							URL:         "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
							Size:        "full",
							AspectRatio: "20:13",
							AspectMode:  "cover",
						},
						Body: &linebot.BoxComponent{
							Type:   "box",
							Layout: "vertical",
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:   "text",
									Text:   "正解です！",
									Weight: "bold",
									Size:   "xl",
								},
								&linebot.BoxComponent{
									Type:     "box",
									Layout:   "vertical",
									Margin:   "md",
									Spacing:  "sm",
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type:   "text",
											Text:   "作者：〇〇",
											Wrap:   true,
											Color:  "#666666",
											Size:   "sm",
										},
									},
								},
							},
						},
						Footer: &linebot.BoxComponent{
							Type:    "box",
							Layout:  "vertical",
							Spacing: "sm",
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type:  "button",
									Style: "link",
									Height: "sm",
									Action: &linebot.URIAction{
										Label: "次の問題",
										URI:   "https://linecorp.com",
									},
								},
								&linebot.ButtonComponent{
									Type:  "button",
									Style: "link",
									Height: "sm",
									Action: &linebot.URIAction{
										Label: "クイズをやめる",
										URI:   "https://linecorp.com",
									},
								},
								&linebot.BoxComponent{
									Type: "box",
									Layout: "vertical",
									Contents: []linebot.FlexComponent{},
									Margin: "sm",
								},
							},
						},
					},
				)
				_, err = bot.ReplyMessage(event.ReplyToken, resp).Do()
				if err != nil {
					log.Print(err)
				}
				}
			}
		// }	else if event.Type == linebot.EventTypePostback {
		// 	// ボタンが押されたときの処理
		// 	if event.Postback.Data == "action=show_another_flex" {
		// 		// 別のFlexメッセージを作成
		// 		anotherFlexMessage := createAnotherFlexMessage()

		// 		// 別のFlexメッセージを送信
		// 		if _, err := bot.ReplyMessage(event.ReplyToken, anotherFlexMessage).Do(); err != nil {
		// 			log.Println("メッセージの送信に失敗しました:", err)
		// 		}
		// 	}
		}
	}
}