package goemail

import (
	"fmt"
	"testing"
	"time"
)

// TODO: 请填写配置信息
func init() {
	c := &Config{
		User:     "xxxx@163.com",
		Password: "xxxx",
		Host:     "smtp.163.com",
		Port:     "465",
		SSL:      true,
		Size:     10,
	}
	SetConfig(c)
	Start()
}

// TODO: 请填写收信人地址
const TO = "xxxx@163.com"

var testIndex int = -1

func TestInit(t *testing.T) {
}

func newTestEmail() *Email {
	testIndex++
	//return NewNormalOneEmail(TO, fmt.Sprintf("Subject(标题)(%v)", testIndex), "This is Content(测试中文内容)", TYPE_PLAIN, errFunc)
	return NewNormalOneEmail(TO, fmt.Sprintf("Subject(标题)(%v)", testIndex), "Hello <b>Bold</b> and <i>Italics</i>!", TYPE_HTML, errFunc)
}

func errFunc(err error) bool {
	panic(err)
	return true
}

func TestOneEmail(t *testing.T) {
	e := newTestEmail()
	if err := SendEmail(e); err != nil {
		t.Error(err)
	}
}

func TestMultiEmail(t *testing.T) {
	e := newTestEmail().SetTo(TO, TO, TO, TO, TO)
	if err := SendEmail(e); err != nil {
		t.Error(err)
	}
}

func TestMultiEmails(t *testing.T) {
	es := []*Email{newTestEmail(), newTestEmail(), newTestEmail()}
	for _, e := range es {
		if err := SendEmail(e); err != nil {
			t.Error(err)
		}
	}
}

func TestWait(t *testing.T) {
	time.Sleep(time.Second)
}
