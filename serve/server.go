package serve

import (
	"context"
	_ "embed"
	"errors"
	"github.com/vaiktorg/grimoire/log"
	"github.com/vaiktorg/grimoire/serve/simws"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	mu sync.Mutex

	AppName string
	addr    string
	tlscfg  *TLSConfig

	handler http.Handler
	ws      simws.ISimWebSocket

	server *http.Server
	Logger log.ILogger
}

type IServer interface {
	ListenAndServe()
	ListenAndServeTLS()
	Startup(init func(s AppConfig))
}

type AppConfig interface {
	WebSocket(wsh func(socket simws.ISimWebSocket))
	MUX(m func(mux *http.ServeMux))
}

func NewServer(config *Config) *Server {
	if config == nil {
		config = defaultConfig
	}

	return &Server{
		AppName: config.GetAppName(),
		handler: config.GetHandler(),
		addr:    config.GetAddr(),
		Logger:  config.GetLoggerConfig(),
		tlscfg:  config.GetTLSConfig(),
		ws:      simws.NewSimWebSocket(context.Background()),
	}
}

func (s *Server) ListenAndServe() {
	defer s.Logger.Close()

	s.server = &http.Server{
		Handler:      s.handler,
		Addr:         s.addr,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}

	// Start the server in a goroutine
	go func() {
		s.Logger.TRACE("Listening " + s.AppName + " on " + s.addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.ERROR("ListenAndServe error: " + err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Block until a signal is received
	<-quit
	s.Logger.TRACE("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	if err := s.server.Shutdown(ctx); err != nil {
		s.Logger.FATAL("ListenAndServe shutdown error: " + err.Error())
	}
}
func (s *Server) ListenAndServeTLS() {
	if s.tlscfg == nil ||
		s.tlscfg.CertPath == "" ||
		s.tlscfg.KeyPath == "" {
		panic("validate tls server config file paths")
	}

	// Define the TLS config in a Server
	s.server = &http.Server{
		Handler:      s.handler,
		Addr:         s.addr,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}
	// Close logger
	defer s.Logger.Close()

	// Start the server in a goroutine
	go func() {
		if err := s.server.ListenAndServeTLS(s.tlscfg.CertPath, s.tlscfg.KeyPath); err != nil || !errors.Is(err, http.ErrServerClosed) {
			s.Logger.FATAL("ListenAndServe error: " + err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Block until a signal is received
	<-quit
	s.Logger.TRACE("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	if err := s.server.Shutdown(ctx); err != nil {
		s.Logger.FATAL("ListenAndServe Shutdown: " + err.Error())
	}

	// Bye bye!
	s.Logger.TRACE(s.AppName + " exiting...")
}

func (s *Server) Startup(init func(cfg AppConfig)) {
	init(s)
}
func (s *Server) WebSocket(wsh func(socket simws.ISimWebSocket)) {
	if sm, ok := s.handler.(*http.ServeMux); ok {
		wsh(s.ws)
		sm.Handle("/ws", s.ws)
	} else {
		panic("server Handler is not of type *http.ServeMux")
	}
}
func (s *Server) MUX(mh func(*http.ServeMux)) {
	if sm, ok := s.handler.(*http.ServeMux); !ok {
		mh(sm)
	} else {
		panic("server Handler is not of type *http.ServeMux")
	}
}
