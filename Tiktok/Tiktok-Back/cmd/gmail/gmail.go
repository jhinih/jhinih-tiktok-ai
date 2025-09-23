package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	backendURL   = "http://localhost:8080/mail"
	pollInterval = 30 * time.Second
)

type mailPayload struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
}

func main() {
	ts := newTokenSource() // 自动刷新 OAuth2 令牌
	for {
		if err := checkGmail(ts); err != nil {
			log.Printf("轮询出错: %v", err)
		}
		time.Sleep(pollInterval)
	}
}

// 创建 OAuth2 TokenSource
func newTokenSource() oauth2.TokenSource {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json: %v", err)
	}
	config, err := google.ConfigFromJSON(b, "https://mail.google.com/")
	if err != nil {
		log.Fatalf("Unable to parse client secret: %v", err)
	}
	tok := &oauth2.Token{}
	tokFile, err := os.Open("token.json")
	if err != nil {
		log.Fatalf("token.json 不存在，请先运行授权工具: %v", err)
	}
	if err := json.NewDecoder(tokFile).Decode(tok); err != nil {
		log.Fatal(err)
	}
	return config.TokenSource(context.Background(), tok)
}

// 连接 Gmail IMAP + 轮询未读
func checkGmail(ts oauth2.TokenSource) error {
	// 1. 拿实时访问令牌
	token, err := ts.Token()
	if err != nil {
		return fmt.Errorf("获取访问令牌失败: %w", err)
	}
	authBytes := []byte("user=" + token.Email + "\x01auth=Bearer " + token.AccessToken + "\x01\x01")

	// 2. 建立 TLS 连接
	addr := "imap.gmail.com:993"
	tlsConfig := &tls.Config{ServerName: "imap.gmail.com"}
	c, err := imapclient.DialTLS(addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接 IMAP 失败: %w", err)
	}
	defer c.Close()

	// 3. OAuth2 认证
	if err := c.Authenticate("XOAUTH2", authBytes).Wait(); err != nil {
		return fmt.Errorf("OAuth2 认证失败: %w", err)
	}
	defer c.Logout().Wait()

	// 4. 选 INBOX
	mbox, err := c.Select("INBOX", false).Wait()
	if err != nil {
		return fmt.Errorf("选择邮箱失败: %w", err)
	}

	// 5. 搜索未读
	criteria := &imap.SearchCriteria{
		Flag: []imap.SearchCriteriaFlag{imap.SearchCriteriaFlagUnseen},
	}
	uids, err := c.Search(criteria, nil).Wait()
	if err != nil {
		return fmt.Errorf("搜索未读失败: %w", err)
	}
	if len(uids) == 0 {
		log.Println("暂无新邮件")
		return nil
	}

	// 6. 取第一封新邮件的头部
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids[0])
	fetchItems := []imap.FetchItem{imap.FetchItemEnvelope}
	msgs, err := c.Fetch(seqSet, fetchItems, nil).Collect()
	if err != nil {
		return fmt.Errorf("抓取邮件失败: %w", err)
	}
	if len(msgs) == 0 {
		return nil
	}
	env := msgs[0].Envelope

	// 7. 推送后端
	p := mailPayload{
		From:    env.From[0].Address(),
		Subject: env.Subject,
		Date:    env.Date.Format(time.RFC3339),
	}
	body, _ := json.Marshal(p)
	resp, err := http.Post(backendURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("推送后端失败: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("已推送新邮件：%s", string(body))
	return nil
}
