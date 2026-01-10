package mockhttp

import (
	"fmt"
	"net/http"
	"time"
)

type MockHttp struct {
	Port    int
	started bool
	// mux     *http.ServeMux
	server *http.Server
}

// Start HTTP server
func (m *MockHttp) Start() {
	if m.Port == 0 {
		panic("no port for mockserver")
	}
	if m.started {
		panic("mockserver already started")
	}

	m.started = true

	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/notfound", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	m.server = &http.Server{Addr: fmt.Sprintf(":%v", m.Port), Handler: mux}

	go func() {
		err := m.server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic("could not listen for mock server")
		}
	}()

	// Readiness check
	counter := 0
	for {
		time.Sleep(time.Millisecond)
		_, err := http.Get(fmt.Sprintf("http://localhost:%v/ok", m.Port))
		if err == nil {
			break
		}

		counter++
		if counter > 20 {
			panic("mock server is not getting ready in time")
		}
	}
}

func (m *MockHttp) Stop() {
	if m.server == nil {
		panic("cannot stop mock server, not started")
	}

	m.server.Close()
}
