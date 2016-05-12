package goemail

import (
	"gopkg.in/gomail.v2"
	"time"
)

type Pool struct {
	pipe   chan *Email
	config *Config
}

func NewPool(c *Config) *Pool {
	if c == nil {
		return new(Pool)
	} else if c.Size < 0 {
		c.Size = 0
	}
	return &Pool{
		pipe:   make(chan *Email, c.Size),
		config: c,
	}
}

func (p *Pool) Start() {
	if p.config == nil || !p.config.Valid() {
		return
	}
	go func() {
		d := gomail.NewDialer(
			p.config.Host,
			p.config.Port,
			p.config.User,
			p.config.Password,
		)
		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case e, ok := <-p.pipe:
				if !ok {
					return
				} else if e == nil || !e.Valid() {
					continue
				} else if !open {
					if s, err = d.Dial(); err != nil {
						if e.FailToSend(err) {
							SendEmail(e)
						}
						continue
					}
					open = true
				}
				if err := gomail.Send(s, e.Message); err != nil {
					if e.FailToSend(err) {
						SendEmail(e)
					}
				} else {
					e.SuccessToSend()
				}
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						//panic(err)
					}
					open = false
				}
			}
		}
	}()
}

func (p *Pool) Stop() {
	close(p.pipe)
}

func (p *Pool) Send(e *Email) error {
	if e == nil || !e.Valid() {
		return ERR_EMAIL
	}
	p.pipe <- e
	return nil
}
