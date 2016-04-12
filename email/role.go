package email

type ROLE string

const ROLE_DEFAULT ROLE = "NORMAL"

func (s ROLE) String() string {
	return string(s)
}
