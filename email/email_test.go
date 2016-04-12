package email

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// TODO: 请填写配置信息
func init() {
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

// TODO: 请填写收信人地址
const TO = "xxxx@163.com"

var testIndex int = -1

func TestInit(t *testing.T) {
}

func makeTestEmail() *Email {
	testIndex++
	return NewOneEmail(ROLE_DEFAULT, TO, fmt.Sprintf("Subject(标题)(%v)", testIndex), "This is Content(测试中文内容)", TYPE_PLAIN)
}

func printTestEmail(t *testing.T, e *Email) {
	t.Logf("%v(%v)#E", e.Subject, e.Tag)
}

func TestOneEmail(t *testing.T) {
	e := makeTestEmail()
	err := SendEmail(e)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiEmail(t *testing.T) {
	e := makeTestEmail()
	e.To = []string{TO, TO, TO, TO, TO}
	err := SendEmail(e)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiEmails(t *testing.T) {
	es := []*Email{makeTestEmail(), makeTestEmail(), makeTestEmail()}
	for _, e := range es {
		err := SendEmail(e)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCheckFailEmail(t *testing.T) {
	for !IsIdle() {
		t.Log("Waited one second...")
		time.Sleep(time.Second)
	}
	for _, e := range GetSuccessEmail() {
		printTestEmail(t, e)
	}
	if HasFailEmail() {
		t.Error(errors.New("Existed unsuccessful email"))
		for _, e := range GetFailEmail() {
			t.Errorf("%v -- %v", e.email.Tag, e.err)
		}
	}
}
