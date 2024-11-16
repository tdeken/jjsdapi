package progress

import (
	"context"
	"fmt"
	grds "github.com/go-redis/redis/v8"
	"jjsdapi/plugins/cache/derivative"
	"strconv"
	"time"
)

const (
	StatusWait    = iota //等待运行
	StatusRunning        //运行中
	StatusSuccess        //导出成功
	StatusFail           //导出失败
)

const (
	defaultExp = 24 * time.Hour //默认过期时间
)

const (
	StatusKey   = "status"   //进度条状态键
	ScheduleKey = "schedule" //当前进度键
	ContentKey  = "content"  //导出后内容键
	ErrMsgKey   = "err_msg"  //导出失败错误信息系键
)

type Progress struct {
	derivative.Derivative
	exp time.Duration
}

type Info struct {
	Status  int     `json:"status"`  //当前状态
	Np      float64 `json:"np"`      //当前进度
	Content string  `json:"content"` //完成内容（进度完成后存储的内容，如：导出完成后的下载地址）
	ErrMsg  string  `json:"err_msg"` //错误信息
}

// NewRedisProgress 实例化redis进度条
func NewRedisProgress(ctx context.Context, cli *grds.Client, key string, opts ...Option) *Progress {
	var pro = &Progress{
		exp: defaultExp,
	}

	for _, opt := range opts {
		opt(pro)
	}

	//初始化
	pro.Init(ctx, cli, key)

	return pro
}

// SetExportStatusRun 设置导出的运行中
func (p *Progress) SetExportStatusRun() {
	p.Cli.HSet(p.Ctx, p.Key, StatusKey, StatusRunning)
}

// SetExportStatusFail 设置导出异常
func (p *Progress) SetExportStatusFail(err error) {
	p.Cli.HSet(p.Ctx, p.Key, StatusKey, StatusFail, ErrMsgKey, err.Error())
}

// SetExportProgress 设置当前导出的进度
func (p *Progress) SetExportProgress(schedule float64) {
	p.Cli.HSet(p.Ctx, p.Key, ScheduleKey, schedule)
}

// SetExportUrl 设置导出的下载地址
func (p *Progress) SetExportUrl(content string) {
	p.Cli.HSet(p.Ctx, p.Key, ContentKey, content, StatusKey, StatusSuccess, ScheduleKey, 100)
}

// GetExportProgress 获取导出进度
func (p *Progress) GetExportProgress() (info Info) {
	res := p.Cli.HMGet(p.Ctx, p.Key, StatusKey, ScheduleKey, ContentKey, ErrMsgKey).Val()

	info.Status, _ = strconv.Atoi(fmt.Sprintf("%v", res[0]))
	info.Np, _ = strconv.ParseFloat(fmt.Sprintf("%v", res[1]), 64)
	info.Content = fmt.Sprintf("%v", res[2])
	info.ErrMsg = fmt.Sprintf("%v", res[3])

	return
}

// Release 清楚导出进度
func (p *Progress) Release() {
	p.Cli.Del(p.Ctx, p.Key)
}

// CreateProgress 创建一个导出进度
func (p *Progress) CreateProgress() {
	p.Cli.HSet(p.Ctx, p.Key, StatusKey, StatusWait, ScheduleKey, 0, ContentKey, "", ErrMsgKey, "")
	p.Cli.Expire(p.Ctx, p.Key, p.exp)
}

// Exist 设置进度条是否存在
func (p *Progress) Exist() bool {
	return p.Cli.Exists(p.Ctx, p.Key).Val() > 0
}
