package email

import "time"

type Config struct {
	User      string
	Password  string
	Host      string
	Port      string
	SSL       bool
	Sleep     time.Duration
	WorkerCnt int
}
