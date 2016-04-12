package email

import "sync"

// TODO: 待被队列替换
type EmailQueue struct {
	queue    []*Email
	queueCnt int
	queueIdx int
	current  *Email
	state    *sync.Mutex
	mutex    *sync.Mutex
}

func NewEmailQueue() *EmailQueue {
	return &EmailQueue{queue: make([]*Email, 0), queueCnt: 0, queueIdx: 0, state: &sync.Mutex{}, mutex: &sync.Mutex{}}
}

func (s *EmailQueue) addEmail(e *Email) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if e == nil {
		return ERR_EMAIL_INVALID
	}
	s.queue = append(s.queue, e)
	s.queueCnt++
	return nil
}

func (s *EmailQueue) removeEmail() (*Email, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if (s.queueCnt - s.queueIdx) <= 0 {
		return nil, ERR_QUEUE_EMPTY
	}
	e := s.queue[s.queueIdx]
	s.queueIdx++
	return e, nil
}

func (s *EmailQueue) IsEmpty() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if (s.queueCnt - s.queueIdx) <= 0 {
		return s.current == nil
	}
	return false
}

func (s *EmailQueue) DispatchToSend(sender baseSender) {
	s.state.Lock()
	defer s.state.Unlock()
	var (
		e   *Email
		err error
	)
	if s.current != nil {
		e = s.current
	} else {
		e, err = s.removeEmail()
		if err != nil {
			return
		}
	}
	err = sender.Send(e)
	if err != nil {
		s.current = e
	} else {
		s.current = nil
	}
}
