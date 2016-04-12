package email

import "errors"

var (
	ERR_EMAIL_INVALID  error = errors.New("Invalid email")
	ERR_EMAIL_ROLE           = errors.New("Invalid role of email")
	ERR_SENDER_BUSY          = errors.New("Sender is busy")
	ERR_WORKER_RUNNING       = errors.New("Worker is rnning")
	ERR_QUEUE_EMPTY          = errors.New("Queue is empty")
)
