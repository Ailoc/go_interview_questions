package main

import (
	"context"
	"net/http"
	"sync"
)

// 设计一个优雅关闭的http服务器

type Server struct {
	srv *http.Server
	wg  sync.WaitGroup
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{Addr: addr, Handler: handler},
		wg:  sync.WaitGroup{},
	}
}

func (s *Server) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
