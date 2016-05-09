package goemail

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	SSL      bool
	Size     int
}

func (c *Config) Valid() bool {
	return len(c.Host) != 0 && c.Port > 0 && len(c.User) != 0 && len(c.Password) != 0
}
