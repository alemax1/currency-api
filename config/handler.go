package config

import "time"

type Handler struct {
	RequestTimeout time.Duration
}

func newHandler() Handler {
	return Handler{
		RequestTimeout: getDefaultDurationEnv("HANDLER_REQUEST_TIMEOUT", 100*time.Millisecond),
	}
}
