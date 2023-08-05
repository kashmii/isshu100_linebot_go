package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

func CorrectMessage(w Waka) linebot.FlexContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:        "image",
			URL:         w.Image,
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
							Text:   "作者：" + w.Author,
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
					Action: linebot.NewPostbackAction("クイズをやめる", "bye", "", ""),
				},
			},
		},
	}
}
