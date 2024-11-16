package webhook

type WebHook interface {
	Write(p []byte) (n int, err error)
	GenSign() (t int64, sign string)
	SendText(content string)
}
