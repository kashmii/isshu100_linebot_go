package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
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
	http.HandleFunc("/callback", lineHandler)

	fmt.Println("http://localhost:8080 で起動中...")
	// HTTPサーバを起動
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func choiceComponent(c Choice) *linebot.ButtonComponent {
	return &linebot.ButtonComponent{
		Type: linebot.FlexComponentTypeButton,
		Style: linebot.FlexButtonStyleTypeLink,
		// PostbackAction(label, data, text, displayText)
		Action: linebot.NewPostbackAction(c.Simo, c.Correct, "", ""),
	}
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
	file, err := os.Open("master.json")
	if err != nil {
		fmt.Println("JSONファイルをオープンできませんでした:", err)
		return
	}
	defer file.Close()

	// ファイル内容をバイト配列に読み込む
	data, err := io.ReadAll(file)
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
			// 設問メッセージ
			case *linebot.TextMessage:
				switch message.Text {
				case "スタート", "次へ", "asd":

					// 和歌を4つ取得して配列：長さ3に入れる
					// 配列[0]は正解として、それ以外ははずれ選択肢として使う
					randomElements := GetRandomWakas(masterData, 4)
					aNum := ""
					choices := make([]Choice, 4)
					for i, ele := range randomElements {
						choices[i].Kami = ele.Kami
						choices[i].Simo = ele.Simo
						if i == 0 {
							aNum = fmt.Sprint(ele.No)
							choices[i].Correct = aNum
						} else {
							choices[i].Correct = "false" +aNum
						}
					}

					question := choices[0].Kami

					// 選択肢4つをランダムに並べ直し
					// シードを設定
					source := rand.NewSource(time.Now().UnixNano())
					rand.New(source)
					rand.Shuffle(len(choices), func(i, j int) {
						choices[i], choices[j] = choices[j], choices[i]
					})
					
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
										Wrap: true,
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
									choiceComponent(choices[0]),
									choiceComponent(choices[1]),
									choiceComponent(choices[2]),
									choiceComponent(choices[3]),
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
		}	else if event.Type == linebot.EventTypePostback {
			answerTypes := []string{`^(?:[1-9][0-9]?|100)$`, `^false(?:[1-9][0-9]?|100)$`}
			// *regexp.Regexp は、regexp パッケージに含まれる Regexp 型へのポインタ
			var regexps []*regexp.Regexp

			for _, pattern := range answerTypes {
				regExp, err := regexp.Compile(pattern)
				if err != nil {
					fmt.Println("正規表現のコンパイルエラー:", err)
					return
				}
				regexps = append(regexps, regExp)
			}
			
			// 最終的に消す
			fmt.Println("postbackのdata: " + event.Postback.Data)
			
			answerStr := event.Postback.Data
			// 正解
			if regexps[0].MatchString(answerStr) {
				fmt.Println(answerStr)

				answerNum, err := strconv.Atoi(answerStr)
				if err != nil {
					fmt.Println("変換エラー:", err)
					return
				}
				answerWaka := masterData[answerNum - 1]
				resp := linebot.NewFlexMessage(
					"this is a flex message",
					CorrectMessage(answerWaka),
				)
				if _, err = bot.ReplyMessage(event.ReplyToken, resp).Do(); err != nil {
					log.Print(err)
				}
			// はずれ
			} else if regexps[1].MatchString(answerStr) {
				answerNum, err := strconv.Atoi(answerStr[5:])
				if err != nil {
					fmt.Println("変換エラー:", err)
					return
				}
				answerWaka := masterData[answerNum - 1]

				resp := linebot.NewFlexMessage(
					"this is a flex message",
					FalseMessage(answerWaka),
				)
				if _, err = bot.ReplyMessage(event.ReplyToken, resp).Do(); err != nil {
					log.Print(err)
				}
			} else if answerStr == "bye" {
				fmt.Println("byebye")
				resp := linebot.NewFlexMessage(
					"this is a flex message",
					QuitPlaying(),
				)
				if _, err = bot.ReplyMessage(event.ReplyToken, resp).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}