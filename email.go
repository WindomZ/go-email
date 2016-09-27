package goemail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
)

const (
	TYPE_HTML  = "html"
	TYPE_PLAIN = "plain"
)

type EmailErrorFunc func(*Email, error) bool

type Email struct {
	Tag       string
	Message   *gomail.Message
	TryCount  int
	ErrorFunc EmailErrorFunc
}

var idxEmail int64 = -1

// NewEmail create a new email object
func NewEmail(m *gomail.Message, fs ...EmailErrorFunc) *Email {
	idxEmail++
	e := &Email{
		Tag:      fmt.Sprintf("e(%v)", idxEmail),
		Message:  m,
		TryCount: 0,
	}
	if fs != nil && len(fs) != 0 {
		e.ErrorFunc = fs[0]
	}
	return e
}

// NewNormalEmail create a normal email object, defind subject, content, type(html or plain)
func NewNormalEmail(subject, content, _type string, fs ...EmailErrorFunc) *Email {
	m := gomail.NewMessage()
	m.SetHeader("Subject", subject)
	switch _type {
	case TYPE_HTML:
		m.SetBody("text/html; charset=UTF-8", content)
	default:
		m.SetBody("text/plain; charset=UTF-8", content)
	}
	return NewEmail(m, fs...)
}

// NewNormalOneEmail create a normal email object, defind subject, content, type(html or plain) and one receiver
func NewNormalOneEmail(to string, subject, content, _type string, fs ...EmailErrorFunc) *Email {
	return NewNormalEmail(subject, content, _type, fs...).SetTo(to)
}

func (e *Email) String() string {
	if to := e.Message.GetHeader("To"); to != nil && len(to) != 0 {
		return fmt.Sprintf("%v:%v", e.Tag, strings.Join(to, " & "))
	}
	return e.Tag
}

func (e *Email) AddTag(tag string) *Email {
	e.Tag += "-" + tag
	return e
}

func (e *Email) SetFrom(add, name string) *Email {
	if len(name) == 0 {
		name = add
	}
	e.Message.SetAddressHeader("From", add, name)
	return e
}

func (e *Email) SetTo(to ...string) *Email {
	e.Message.SetHeader("To", to...)
	return e
}

func (e *Email) Increase() *Email {
	e.TryCount++
	return e
}

func (e *Email) Valid() bool {
	return e.Message != nil
}

func (e *Email) SuccessToSend() {
	if e != nil && e.ErrorFunc != nil {
		e.ErrorFunc(e, nil)
	}
}

func (e *Email) FailToSend(err error) bool {
	if e != nil && e.ErrorFunc != nil {
		return e.ErrorFunc(e, err)
	}
	return false
}
