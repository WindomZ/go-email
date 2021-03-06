package goemail

import (
	"fmt"
	"testing"
	"time"
)

// TODO: 请填写配置信息
var config = &Config{
	Host:     "smtp.163.com",
	Port:     465,
	User:     "xxxx@163.com",
	Password: "xxxx",
	SSL:      true,
	Size:     10,
}

// TODO: 请填写收信人地址
const TO = "xxxx@163.com"

var testIndex int = -1

func newTestEmail() *Email {
	testIndex++
	return NewNormalOneEmail(
		TO,
		fmt.Sprintf("Subject(标题)(%v)", testIndex),
		"This is Content(测试中文内容)",
		TYPE_PLAIN,
		panicFunc,
	)
	//return NewNormalOneEmail(
	//	TO,
	//	fmt.Sprintf("Subject(标题)(%v)", testIndex),
	//	"Hello <b>Bold</b> and <i>Italics</i>!",
	//	TYPE_HTML,
	//	panicFunc,
	//)
}

func panicFunc(e *Email, err error) bool {
	if err != nil {
		panic(err)
	}
	return true
}

func TestInit(t *testing.T) {
	SetConfig(config)
	StartService()
}

func TestOneEmail(t *testing.T) {
	e := newTestEmail()
	if err := SendEmail(e); err != nil {
		t.Fatal(err)
	}
}

func TestMultiEmail(t *testing.T) {
	e := newTestEmail().SetTo(TO, TO, TO, TO, TO)
	if err := SendEmail(e); err != nil {
		t.Fatal(err)
	}
}

func TestMultiEmails(t *testing.T) {
	es := []*Email{newTestEmail(), newTestEmail(), newTestEmail()}
	for _, e := range es {
		if err := SendEmail(e); err != nil {
			t.Fatal(err)
		}
	}
}

func TestWait(t *testing.T) {
	time.Sleep(time.Second * 5)
	StopService()
}
