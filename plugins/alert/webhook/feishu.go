package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	feiShuChanLen = 100
)

type FeiShu struct {
	Url      string
	Secret   string
	sendChan chan string
}

func NewFeiShu(url string, secret string) *FeiShu {
	f := FeiShu{
		Url:      url,
		Secret:   secret,
		sendChan: make(chan string, feiShuChanLen),
	}
	f.MakeSendTask()
	return &f
}

func (f FeiShu) GenSign() (t int64, sign string) {
	t = time.Now().Unix()
	stringToSign := fmt.Sprint(t) + "\n" + f.Secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	h.Write(data)
	sign = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

func (f FeiShu) SendText(content string) {
	if len(f.sendChan) < feiShuChanLen {
		f.sendChan <- content
	}
}

func (f FeiShu) Write(p []byte) (n int, err error) {
	f.SendText(string(p))
	return
}

func (f FeiShu) MakeSendTask() {
	go func() {
		defer func() {
			if re := recover(); re != nil {
				time.Sleep(5 * time.Second)
				f.MakeSendTask()
				return
			}
		}()
		for content := range f.sendChan {
			f.SendTask(content)
		}
	}()
}

func (f FeiShu) SendTask(content string) {
	var data = map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": content,
		},
	}
	if f.Secret != "" {
		t, sign := f.GenSign()
		data["timestamp"] = t
		data["sign"] = sign
	}

	d, _ := json.Marshal(data)
	c := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", f.Url, bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/json")
	c.Do(req)
}
