package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Server *http.Server
}

func New(h http.Handler, port int) *Server {
	addr := fmt.Sprintf(":%d", port)
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: h,
		},
	}
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
