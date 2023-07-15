package app

type Server struct {
	Host string
	Port int64
}

// NewServer creates a new server with host and port attributes.
func NewServer(host string, port int64) *Server {
	return &Server{
		Port: port,
		Host: host,
	}
}
