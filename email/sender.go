package email

import (
	"fmt"
	"sync"
)

type baseSender interface {
	IsIdle() bool
	Send(e *Email) error
	fetch()
	pushSuccess(*Email)
	pushFailure(*Email, error)
}

type Sender struct {
	baseSender
	Tag        string
	dispatcher baseDispatcher
	config     *Config
	workers    []*Worker
	mutex      *sync.Mutex
}

func NewSender(tag string, d baseDispatcher, c *Config) *Sender {
	s := &Sender{Tag: tag, dispatcher: d, config: c, mutex: &sync.Mutex{}}
	s.workers = make([]*Worker, c.WorkerCnt)
	for i := 0; i < c.WorkerCnt; i++ {
		s.workers[i] = NewWorker(fmt.Sprintf("wk(%v)", i), c, s)
	}
	return s
}

func (s *Sender) IsIdle() bool {
	for _, w := range s.workers {
		if !w.IsIdle() {
			return false
		}
	}
	return true
}

func (s *Sender) Send(e *Email) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if e == nil {
		return ERR_EMAIL_INVALID
	}
	var err error = ERR_SENDER_BUSY
	e.AddTag(s.Tag)
	for _, w := range s.workers {
		if !w.isRunning {
			err = w.Send(e)
			break
		}
	}
	if err != nil && err != ERR_SENDER_BUSY {
		s.pushFailure(e, err)
	}
	return err
}

func (s *Sender) fetch() {
	if !s.IsIdle() {
		return
	}
	s.dispatcher.dispatch()
}

func (s *Sender) pushSuccess(e *Email) {
	appendSuccessEmail(e)
}

func (s *Sender) pushFailure(e *Email, err error) {
	appendFailEmail(e, err)
}

// TODO: 增加定时器
