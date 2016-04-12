package email

import (
	"runtime"
	"time"
)

type failEmail struct {
	email *Email
	err   error
}

var (
	dispatchers   map[string]*Dispatcher
	successEmails []*Email
	failEmails    []*failEmail
)

func init() {
	// TODO: 配置信息加载 -> 初始化
	runtime.GOMAXPROCS(2) // 待定
	c := &Config{User: "xxxx@163.com",
		Password:  "xxxx",
		Host:      "smtp.163.com",
		Port:      "465",
		SSL:       true,
		Sleep:     3000 * time.Millisecond,
		WorkerCnt: 3,
	}
	cs := make([]*Config, 3)
	for i := 0; i < 3; i++ {
		cs[i] = c
	}
	Init(cs)
}

func Init(cs []*Config) {
	dispatchers = make(map[string]*Dispatcher, 1)
	dispatchers[ROLE_DEFAULT.String()] = NewDispatcher(ROLE_DEFAULT, cs)
}

func addEmail(e *Email) error {
	if e == nil {
		return ERR_EMAIL_INVALID
	}
	switch e.Role {
	case ROLE_DEFAULT:
		return dispatchers[ROLE_DEFAULT.String()].AddEmail(e)
	default:
		return dispatchers[ROLE_DEFAULT.String()].AddEmail(e)
	}
}

func SendEmail(e *Email) error {
	return addEmail(e)
}

func IsIdle() bool {
	var r bool = true
	for _, d := range dispatchers {
		r = r && d.IsIdle()
	}
	return r
}

func appendSuccessEmail(e *Email) {
	successEmails = append(successEmails, e)
}

func GetSuccessEmail() []*Email {
	return successEmails
}

func appendFailEmail(e *Email, err error) {
	failEmails = append(failEmails, &failEmail{email: e, err: err})
	// TODO： 处理失败邮件
}

func HasFailEmail() bool {
	return len(failEmails) != 0
}

func GetFailEmail() []*failEmail {
	return failEmails
}
