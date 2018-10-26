package server

import "github.com/aweris/stp/storage"

// Server is a wrapper object for internal services
type Server struct {
	db *storage.BoltDB
}

// NewServer creates and configures with boltDB storage
func NewServer(storagePath string) *Server {
	db, err := storage.NewBoltDB(storagePath)
	if err != nil {
		return nil
	}

	s := &Server{db: db}

	return s
}

func (s *Server) Close() {
	s.db.Close()
}
