package ratelimit

type Config struct {
	Manager Manager
	Force   bool
}

var defaultConfig = Config{
	Manager: nil,
	Force:   false,
}
