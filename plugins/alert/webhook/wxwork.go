package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const (
	wxWorkChanLen = 100
)

type WxWork struct {
	Url      string
	Secret   string
	sendChan chan string
}

func NewWxWork(url string, secret string) *WxWork {
	f := WxWork{
		Url:      url,
		Secret:   secret,
		sendChan: make(chan string, wxWorkChanLen),
	}
	f.MakeSendTask()
	return &f
}

func (w WxWork) GenSign() (t int64, sign string) {
	return
}

func (w WxWork) SendText(content string) {
	if len(w.sendChan) < wxWorkChanLen {
		w.sendChan <- content
	}
}

func (w WxWork) Write(p []byte) (n int, err error) {
	w.SendText(string(p))
	return
}

func (w WxWork) MakeSendTask() {
	go func() {
		defer func() {
			if re := recover(); re != nil {
				time.Sleep(5 * time.Second)
				w.MakeSendTask()
				return
			}
		}()
		for content := range w.sendChan {
			w.SendTask(content)
		}
	}()
}

func (w WxWork) SendTask(content string) {
	var data = map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}
	d, _ := json.Marshal(data)
	c := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", w.Url, bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/json")
	c.Do(req)
}
