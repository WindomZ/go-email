package goemail

import "errors"

var (
	ERR_CONFIG         error = errors.New("goemail: Invalid config")
	ERR_EMAIL                = errors.New("goemail: Invalid email")
	ERR_EMAIL_TOO_MUCH       = errors.New("goemail: Try to send email too much")
	ERR_POOL                 = errors.New("goemail: Invalid pool")
	ERR_FULL_POOL            = errors.New("goemail: full pool")
)
