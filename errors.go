package goemail

import "errors"

var (
	ERR_CONFIG         error = errors.New("goemail: invalid config")
	ERR_EMAIL                = errors.New("goemail: invalid email")
	ERR_EMAIL_TOO_MUCH       = errors.New("goemail: try to send email too much")
	ERR_POOL                 = errors.New("goemail: invalid pool")
	ERR_FULL_POOL            = errors.New("goemail: full pool")
)
