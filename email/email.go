package email

import "fmt"

const (
	TYPE_HTML  = "html"
	TYPE_PLAIN = "plain"
)

type Email struct {
	Tag       string
	Role      ROLE
	To        []string
	Subject   string
	Content   string
	Type      string
	MaxTryCnt int
}

var idxEmail int64 = -1

func NewEmail(r ROLE, to []string, s string, c string, t string) *Email {
	var conTyp string
	switch t {
	case TYPE_HTML:
		conTyp = "Content-Type: text/html; charset=UTF-8"
	default:
		conTyp = "Content-Type: text/plain; charset=UTF-8"
	}
	idxEmail++
	return &Email{Role: r, Tag: fmt.Sprintf("e(%v)", idxEmail), To: to, Subject: s, Content: c, Type: conTyp, MaxTryCnt: 3}
}

func NewOneEmail(r ROLE, to string, s string, c string, t string) *Email {
	return NewEmail(r, []string{to}, s, c, t)
}

func (s *Email) AddTag(tag string) {
	s.Tag += "-" + tag
}
