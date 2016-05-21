package goemail

import (
	"strconv"
	"time"
)

var pools []*Pool

func SetConfig(cs ...*Config) {
	if cs == nil || len(cs) == 0 {
		panic(ERR_CONFIG)
	}
	pools = make([]*Pool, 0, len(cs))
	for _, c := range cs {
		if c == nil || !c.Valid() {
			continue
		}
		pools = append(pools, NewPool(c))
	}
}

func StartService() error {
	if pools == nil || len(pools) == 0 {
		return ERR_POOL
	}
	for _, p := range pools {
		if p != nil {
			p.Start()
		}
	}
	return nil
}

func StopService() error {
	if pools == nil || len(pools) == 0 {
		return ERR_POOL
	}
	for _, p := range pools {
		if p != nil {
			p.Stop()
		}
	}
	return nil
}

var sendIdx int = 0

func SendEmail(e *Email) error {
	if e == nil || !e.Valid() {
		return ERR_EMAIL
	} else if pools == nil || len(pools) == 0 {
		return ERR_POOL
	}
	if sendIdx > 0 && sendIdx >= len(pools) {
		sendIdx = 0
	}
	if p := pools[sendIdx]; p != nil {
		sendIdx++
		if e.TryCount >= 5 {
			return ERR_EMAIL_TOO_MUCH
		} else if err := p.Send(e.
			AddTag(strconv.Itoa(sendIdx)).
			SetFrom(p.config.User).
			Increase()); err != nil {
			return SendEmail(e)
		}
	}
	return nil
}

func SendEmailDelay(e *Email) {
	go func() {
		time.Sleep(time.Second * 3)
		SendEmail(e)
	}()
}
