package email

import (
	"fmt"
	"github.com/WindomZ/go-dice/dice"
	"runtime"
)

type baseDispatcher interface {
	dispatch()
}

type Dispatcher struct {
	Role    ROLE
	queue   *EmailQueue
	senders []*Sender
	dice    *dice.DiceInt
}

func NewDispatcher(r ROLE, cs []*Config) *Dispatcher {
	d := &Dispatcher{Role: r, senders: make([]*Sender, len(cs)), dice: dice.NewDiceInt(len(cs), dice.TYPE_POLL)}
	d.queue = NewEmailQueue()
	for i, c := range cs {
		d.senders[i] = NewSender(fmt.Sprintf("sd(%v)(%v)", r.String(), i), d, c)
	}
	return d
}

func (s *Dispatcher) AddEmail(e *Email) error {
	if e == nil {
		return ERR_EMAIL_INVALID
	}
	if e.Role != s.Role {
		return ERR_EMAIL_ROLE
	}
	err := s.queue.addEmail(e)
	if err == nil {
		s.dispatch()
	}
	return err
}

func (s *Dispatcher) IsIdle() bool {
	var r bool = s.queue.IsEmpty()
	if !r {
		return false
	}
	for _, s := range s.senders {
		r = r && s.IsIdle()
		if !r {
			break
		}
	}
	return r
}

func (s *Dispatcher) dispatch() {
	go s.dispatchGo()
}

func (s *Dispatcher) dispatchGo() {
	runtime.Gosched()
	s.queue.DispatchToSend(s.senders[s.dice.Throw().Value])
}
