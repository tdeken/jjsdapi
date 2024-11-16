package certs

import (
	"certified/internal/utils"
	"encoding/json"
	"github.com/rs/xid"
	"testing"
)

func TestName(t *testing.T) {
	t.Log(xid.New().String())
	for i := 0; i < 10; i++ {
		t.Logf(randomString(10))
	}
}

func TestUser(t *testing.T) {

	var token = "eyJhbGciOiJIUzI1NiIsInNhbHQiOiJ1R1V0MEciLCJzaWduIjoiMjY1YjRhMDBlMTg3MjViY2NkM2QzMTljMDFhY2YxZTUiLCJ0eXAiOiJKV1QiLCJ2ZXIiOiJ2MSJ9.eyJpbnRfaWQiOjE3NDEyLCJjYW1wYWlnbl9pZCI6MjgwMywiY2hhbm5lbF9pZCI6NTQ1OCwicGlkIjoxNzQxMiwiaXNfdmlldyI6dHJ1ZSwibGFuZyI6IiIsImx2IjoiIiwic2hvd19sYW5nIjoiemgiLCJpc3MiOiJjZXJ0aWZpZWQiLCJleHAiOjE3MTQzNzIyNzAsIm5iZiI6MTcxNDM2NTA3MCwiaWF0IjoxNzE0MzY1MDcwLCJqdGkiOiJjb25pNTNtcmRqN3FxOHBhdmV1MCJ9.m_goJaAxE2pqBvxsXM-ijIC7OxggcWqU2ZZCYlMMgHw"

	//t.Run("token", func(t *testing.T) {
	//	user := NewUser()
	//	var err error
	//	token, err = user.Token(1, "username")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	t.Log(token)
	//})

	t.Run("parse", func(t *testing.T) {
		user := NewUser()
		t.Log(user.Parse(token))

		t.Log(user)
	})

}

func TestFUser(t *testing.T) {

	//var token = "eyJhbGciOiJIUzI1NiIsInNhbHQiOiJQMEV6ejYiLCJzaWduIjoiMjdhMTFhODBhZDgwYzg5NGYwMzZlMDI2MDM4NzhkOGQiLCJ0eXAiOiJKV1QiLCJ2ZXIiOiJ2MSJ9.eyJpbnRfaWQiOjE4MjkzOTEyNTAxODI0NzE2OCwiY2FtcGFpZ25faWQiOjM5MzIsImNoYW5uZWxfaWQiOjc1NjYsInBpZCI6MTgyOTM5MTI1MDE4MjQ3MTY4LCJpc192aWV3Ijp0cnVlLCJsYW5nIjoiemgiLCJzaG93X2xhbmciOiJ6aCIsInR6IjoiR01UIiwidHpfb2Zmc2V0Ijo4LCJpc3MiOiJjZXJ0aWZpZWQiLCJleHAiOjE3MzA4NjgxNTcsIm5iZiI6MTczMDg2MDk1NywiaWF0IjoxNzMwODYwOTU3LCJqdGkiOiJjc2xkZjdlcmRqN282dXVzaGtrZyJ9.nfTsi_LN-l6-K30oMAuvJiNUpPg7JvqstgoGgvkBLeE"

	//t.Run("token", func(t *testing.T) {
	//	user := NewUser()
	//	var err error
	//	token, err = user.Token(1, "username")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	t.Log(token)
	//})qian

	t.Run("parse", func(t *testing.T) {
		user := NewFUser()

		tk, _ := user.Token(&FUser{
			BID:        utils.LongNumIdToStr(182939125018247168),
			CampaignId: 3932,
			ChannelId:  7566,
			Pid:        0,
			IsView:     true,
			Lang:       "zh",
			Lv:         "",
			ShowLang:   "zh",
			Country:    "",
			Tz:         "tz",
			TzOffset:   8,
		})

		t.Log(user.Parse(tk))

		t.Log(user)
		t.Log(user.ID, user.FUser.BID)
		t.Log(user.FUser)
		t.Logf("匹配中的语种：%s, 用户进入活动带的语种：%s", user.Lang, user.ShowLang)
	})

}

type TT struct {
	ID int64 `json:"id"`
}

func TestName1(t *testing.T) {
	tt := TT{ID: 182939125018247168}

	b, _ := json.Marshal(tt)
	t.Log(string(b))
}
