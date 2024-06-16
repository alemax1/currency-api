package config

type Server struct {
	Port uint
}

func newServer() Server {
	return Server{
		Port: uint(getDefaultIntEnv("SERVER_PORT", 3000)),
	}
}
