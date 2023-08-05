package main

import (
	"fmt"
	"log"
	"net/http"
	// "github.com/joho/godotenv"
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
	// 本番環境ではこれは使えないぽい
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic("Error loading .env file")
	// }

	// ハンドラの登録
	http.HandleFunc("/callback", LineHandler)

	fmt.Println("http://localhost:8080 で起動中...")
	// HTTPサーバを起動
	log.Fatal(http.ListenAndServe(":8080", nil))
}
