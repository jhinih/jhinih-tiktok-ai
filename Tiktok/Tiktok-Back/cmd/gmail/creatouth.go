package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// 权限只需要 readonly 即可
var config *oauth2.Config

func init() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json: %v", err)
	}
	config, err = google.ConfigFromJSON(b, "https://mail.google.com/")
	if err != nil {
		log.Fatalf("Unable to parse client secret: %v", err)
	}
}

func main() {
	// 生成授权链接
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	println("请访问并复制授权码：\n", url)

	var code string
	print("粘贴授权码：")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	// 换 token
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	f, _ := os.Create("token.json")
	defer f.Close()
	json.NewEncoder(f).Encode(tok)
	println("token.json 已生成，放到项目根目录即可")
}
