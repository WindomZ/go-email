package email

import (
	"crypto/tls"
	"git.yichui.net/zj/mall/api/libs/uuid"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"runtime"
	"sync"
	"time"
)

type Worker struct {
	baseSender
	Tag       string
	Client    *smtp.Client
	config    *Config
	sender    baseSender
	address   *mail.Address
	isValid   bool
	isRunning bool
	tryCnt    int
	mutex     *sync.Mutex
}

func NewWorker(tag string, c *Config, s baseSender) *Worker {
	return &Worker{Tag: tag, config: c, sender: s, address: &mail.Address{c.User, c.User}, isValid: true, isRunning: false, tryCnt: 0, mutex: &sync.Mutex{}}
}

func (s *Worker) initClient() error {
	var (
		auth = smtp.PlainAuth(uuid.NewUUID(), s.config.User, s.config.Password, s.config.Host)
		host = s.config.Host + ":" + s.config.Port
		tc   = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.config.Host,
		}
		cnn net.Conn
		err error
	)
	if s.config.SSL == true {
		if cnn, err = tls.Dial("tcp", host, tc); err != nil {
			return err
		}
		if s.Client, err = smtp.NewClient(cnn, s.config.Host); err != nil {
			return err
		}
	} else {
		if s.Client, err = smtp.Dial(host); err != nil {
			return err
		}
	}
	if err = s.Client.Auth(auth); err != nil {
		return err
	}
	s.isValid = (err == nil)
	return err
}

func (s *Worker) quitClient() error {
	return s.Client.Quit()
}

func (s *Worker) IsIdle() bool {
	return !s.isRunning
}

func (s *Worker) Send(e *Email) error {
	if e == nil {
		return ERR_EMAIL_INVALID
	}
	e.AddTag(s.Tag)
	if !s.IsIdle() {
		return ERR_WORKER_RUNNING
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := s.initClient(); err != nil {
		return err
	}
	s.isRunning = true
	go s.handleEmail(e)
	return nil
}

func (s *Worker) handleEmail(e *Email) {
	runtime.Gosched()
	s.mutex.Lock()
	defer func() {
		s.isRunning = false
		s.mutex.Unlock()
		s.fetch()
	}()
	s.reset()
	for _, t := range e.To {
		s.sendOne(t, e)
	}
	s.quitClient()
}

func (s *Worker) fetch() {
	s.sender.fetch()
}

func (s *Worker) reset() {
	s.tryCnt = 0
}

func (s *Worker) sendOne(t string, e *Email) {
	var (
		err error
		msg []byte
		wc  io.WriteCloser
	)
	msg = []byte("To: " + t + "\r\n" +
		"From: " + s.address.String() + "\r\n" +
		"Subject: " + e.Subject + "\r\n" +
		e.Type + "\r\n\r\n" +
		e.Content)
	err = s.Client.Mail(s.address.Address)
	if err != nil {
		goto final
	}
	err = s.Client.Rcpt(t)
	if err != nil {
		goto final
	}
	wc, err = s.Client.Data()
	if err != nil {
		goto final
	}
	_, err = wc.Write(msg)
	if err != nil {
		goto final
	}
	err = wc.Close()
	if err != nil {
		goto final
	}
final:
	if err != nil {
		s.tryCnt++
		if s.tryCnt <= e.MaxTryCnt {
			time.Sleep(s.config.Sleep)
			s.sendOne(t, e)
		} else {
			s.pushOneFailure(t, *e, err)
		}
	} else {
		s.pushSuccess(e)
	}
}

func (s *Worker) pushSuccess(e *Email) {
	s.sender.pushSuccess(e)
}

func (s *Worker) pushFailure(e *Email, err error) {
	s.sender.pushFailure(e, err)
}

func (s *Worker) pushOneFailure(to string, e Email, err error) {
	e.To = []string{to}
	s.pushFailure(&e, err)
}
