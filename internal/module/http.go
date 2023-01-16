package module

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTP struct {
	Addr string
	DB   *sql.DB
}

func (m *HTTP) Init(ctx context.Context) (err error) {
	m.Addr = ":8080"
	return
}

func (m *HTTP) Run(ctx context.Context) (err error) {
	if m.DB == nil {
		return fmt.Errorf("DB is requried")
	}

	server := http.Server{
		Addr:    m.Addr,
		Handler: m.router(),
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.TODO())
	}()
	log.Printf("Listen at %s", m.Addr)
	if err = server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return
}

func (m *HTTP) router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthz"))
	r.Get("/hello", m.hello)
	return r
}

func (m *HTTP) hello(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintln(writer, "Hello!")
	if err != nil {
		log.Printf("HTTP I/O error: %s", err)
	}
}
