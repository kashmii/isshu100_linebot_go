package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
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

	fmt.Println(len(masterData))
	fmt.Println("=== ===")

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
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										// PostbackAction(label, data, text, displayText)
										Action: linebot.NewPostbackAction(choices[0].Simo, choices[0].Correct, "", ""),
									},
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										Action: linebot.NewPostbackAction(choices[1].Simo, choices[1].Correct, "", ""),
									},
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										Action: linebot.NewPostbackAction(choices[2].Simo, choices[2].Correct, "", ""),
									},
									&linebot.ButtonComponent{
										Type: linebot.FlexComponentTypeButton,
										Style: linebot.FlexButtonStyleTypeLink,
										Action: linebot.NewPostbackAction(choices[3].Simo, choices[3].Correct, "", ""),
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
							URL:         masterData[2].Image,
							Size:        "4xl",
							AspectRatio: "2:3",
							AspectMode:  "cover",
						},
						Body: &linebot.BoxComponent{
							Type:   "box",
							Layout: "vertical",
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:   "text",
									Text:   "正解！",
									Weight: "bold",
									Size:   "xl",
									Align: "center",
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
									Action: linebot.NewMessageAction("次の問題", "flex"),
								},
								&linebot.ButtonComponent{
									Type:  "button",
									Style: "link",
									Height: "sm",
									Action: linebot.NewMessageAction("クイズをやめる", "bye"),
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
		// 正解
		}	else if event.Type == linebot.EventTypePostback {
			answerStr := event.Postback.Data
			// 0~9と100?のとき判別できない
			pattern1 := `^[0-9]{1}$`
			pattern2 := `^[0-9]{2}$`
			pattern3 := `^false[0-9]{1}$`
			pattern4 := `^false[0-9]{2}$`
			// 正規表現のコンパイル
			regExp1, err := regexp.Compile(pattern1)
			if err != nil {
				fmt.Println("正規表現のコンパイルエラー:", err)
				return
			}
			regExp2, err := regexp.Compile(pattern2)
			if err != nil {
				fmt.Println("正規表現のコンパイルエラー:", err)
				return
			}
			regExp3, err := regexp.Compile(pattern3)
			if err != nil {
				fmt.Println("正規表現のコンパイルエラー:", err)
				return
			}
			regExp4, err := regexp.Compile(pattern4)
			if err != nil {
				fmt.Println("正規表現のコンパイルエラー:", err)
				return
			}

			fmt.Println(event.Postback.Data, ":", reflect.TypeOf(event.Postback.Data))


			// はずれ
			if regExp3.MatchString(answerStr) || regExp4.MatchString(answerStr) {
				fmt.Println("JSONファイル")
				answerNum, err := strconv.Atoi(answerStr[5:])
				if err != nil {
					fmt.Println("変換エラー:", err)
					return
				}

				answerWaka := masterData[answerNum - 1]

				resp := linebot.NewFlexMessage(
					"this is a flex message",
					&linebot.BubbleContainer{
						Type: linebot.FlexContainerTypeBubble,
						Body: &linebot.BoxComponent{
							Type:   "box",
							Layout: "vertical",
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:   "text",
									Text:   "残念…",
									Weight: "bold",
									Size:   "xl",
									Align: "center",
								},
								&linebot.BoxComponent{
									Type:     "box",
									Layout:   "vertical",
									Margin:   "md",
									Spacing:  "sm",
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type:   "text",
											Text:   "上の句：" + answerWaka.Kami,
											Wrap:   true,
											Color:  "#666666",
											Size:   "sm",
										},
										&linebot.TextComponent{
											Type:   "text",
											Text:   "下の句：" + answerWaka.Simo,
											Wrap:   true,
											Color:  "#666666",
											Size:   "sm",
										},
										&linebot.TextComponent{
											Type:   "text",
											Text:   "作者：" + answerWaka.Author,
											Wrap:   true,
											Color:  "#666666",
											Size:   "sm",
										},
									},
								},
							},
						},
						Footer: &linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type:  "button",
									Style: "link",
									Height: "sm",
									Action: linebot.NewMessageAction("次の問題", "次へ"),
								},
								&linebot.ButtonComponent{
									Type:  "button",
									Style: "link",
									Height: "sm",
									Action: linebot.NewMessageAction("クイズをやめる", "bye"),
								},
							},
						},
					},
				)
				
				if _, err = bot.ReplyMessage(event.ReplyToken, resp).Do(); err != nil {
					log.Print(err)
				}



			} else if regExp1.MatchString(answerStr) || regExp2.MatchString(answerStr) {
				fmt.Println(len(masterData))
				fmt.Println("=== ===")
				fmt.Println("seikai")

				fmt.Println(answerStr)

				answerNum, err := strconv.Atoi(answerStr)
				if err != nil {
					fmt.Println("変換エラー:", err)
					return
				}
				answerWaka := masterData[answerNum - 1]
				resp := linebot.NewFlexMessage(
					"this is a flex message",
					&linebot.BubbleContainer{
						Type: linebot.FlexContainerTypeBubble,
						Hero: &linebot.ImageComponent{
							Type:        "image",
							URL:         answerWaka.Image,
							Size:        "4xl",
							AspectRatio: "2:3",
							AspectMode:  "cover",
						},
						Body: &linebot.BoxComponent{
							Type:   "box",
							Layout: "vertical",
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:   "text",
									Text:   "正解！",
									Weight: "bold",
									Size:   "xl",
									Align: "center",
								},
								&linebot.BoxComponent{
									Type:     "box",
									Layout:   "vertical",
									Margin:   "md",
									Spacing:  "sm",
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type:   "text",
											Text:   "作者：" + answerWaka.Author,
											Wrap:   true,
											Color:  "#666666",
											Size:   "sm",
										},
									},
								},
							},
						},
						Footer: &linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type:  "button",
									Style: "link",
									Height: "sm",
									Action: linebot.NewMessageAction("次の問題", "次へ"),
								},
								&linebot.ButtonComponent{
									Type:  "button",
									Style: "link",
									Height: "sm",
									Action: linebot.NewMessageAction("クイズをやめる", "bye"),
								},
							},
						},
					},
				)
				
				if _, err = bot.ReplyMessage(event.ReplyToken, resp).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}