package goemail

import (
	"gopkg.in/gomail.v2"
	"time"
)

type Pool struct {
	pipe   chan *Email
	error  chan error
	config *Config
}

func NewPool(c *Config) *Pool {
	if c == nil {
		return new(Pool)
	} else if c.Size < 0 {
		c.Size = 10
	}
	return &Pool{
		pipe:   make(chan *Email, c.Size),
		error:  make(chan error),
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
		var s gomail.SendCloser = nil
		var err error
		for {
			select {
			case e, ok := <-p.pipe:
				if !ok {
					return
				} else if e == nil || !e.Valid() {
					continue
				} else if s == nil {
					if s, err = d.Dial(); err != nil {
						if e.FailToSend(err) {
							SendEmail(e)
						}
						continue
					}
				}
				if err := gomail.Send(s, e.Message); err != nil {
					p.error <- err
					if e.FailToSend(err) {
						SendEmailDelay(e)
					}
				} else {
					e.SuccessToSend()
				}
			case e, ok := <-p.error:
				if !ok {
					return
				} else if e != nil {
					if s != nil {
						if err := s.Close(); err != nil {
						}
						s = nil
					}
				}
			case <-time.After(30 * time.Second):
				if s != nil {
					if err := s.Close(); err != nil {
						//panic(err)
					}
					s = nil
				}
			}
		}
	}()
}

func (p *Pool) Stop() {
	close(p.pipe)
}

func (p *Pool) IsFull() bool {
	return len(p.pipe) >= p.config.Size
}

func (p *Pool) Send(e *Email) error {
	if e == nil || !e.Valid() {
		return ERR_EMAIL
	} else if p.IsFull() {
		return ERR_FULL_POOL
	}
	p.pipe <- e
	return nil
}
