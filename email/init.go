package email

type failEmail struct {
	email *Email
	err   error
}

var (
	dispatchers   map[string]*Dispatcher
	successEmails []*Email
	failEmails    []*failEmail
)

// TODO: 配置信息加载 -> 初始化
func Init(cs []*Config) {
	dispatchers = make(map[string]*Dispatcher, 1)
	dispatchers[ROLE_DEFAULT.String()] = NewDispatcher(ROLE_DEFAULT, cs)
}

func addEmail(e *Email) error {
	if e == nil {
		return ERR_EMAIL_INVALID
	}
	switch e.Role {
	case ROLE_DEFAULT:
		return dispatchers[ROLE_DEFAULT.String()].AddEmail(e)
	default:
		return dispatchers[ROLE_DEFAULT.String()].AddEmail(e)
	}
}

func SendEmail(e *Email) error {
	return addEmail(e)
}

func IsIdle() bool {
	var r bool = true
	for _, d := range dispatchers {
		r = r && d.IsIdle()
	}
	return r
}

func appendSuccessEmail(e *Email) {
	successEmails = append(successEmails, e)
}

func GetSuccessEmail() []*Email {
	return successEmails
}

func appendFailEmail(e *Email, err error) {
	failEmails = append(failEmails, &failEmail{email: e, err: err})
	// TODO： 处理失败邮件
}

func HasFailEmail() bool {
	return len(failEmails) != 0
}

func GetFailEmail() []*failEmail {
	return failEmails
}
