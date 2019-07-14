package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	hs, logger := setup()
	go func() {
		logger.Printf("Listening on http://localhost%s\n", hs.Addr)

		if err := hs.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()
	graceful(hs, logger)
}

func setup() (*http.Server, *log.Logger) {
	port := ":8010"
	hs := &http.Server{Addr: port, Handler: &server{}}

	return hs, log.New(os.Stdout, "", 0)
}

func graceful(hs *http.Server, logger *log.Logger) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Println("\nServer gracefully stopped")
}
