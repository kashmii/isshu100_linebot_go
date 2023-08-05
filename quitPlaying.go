package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

func QuitPlaying() linebot.FlexContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   "box",
			Layout: "vertical",
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:     "box",
					Layout:   "vertical",
					Margin:   "md",
					Spacing:  "sm",
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   "text",
							Text:   "また始めるときは下のボタンを押してくださいね！",
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
					Action: linebot.NewMessageAction("問題を始める", "スタート"),
				},
			},
		},
	}
}
