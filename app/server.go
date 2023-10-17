package app

type Server struct {
	Host            string
	Port            int
	ShutdownTimeout int
}

// NewServer creates a new server with host and port attributes.
func NewServer(host string, port int) *Server {
	return &Server{
		Port: port,
		Host: host,
	}
}
