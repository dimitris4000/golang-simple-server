package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	version                = "0.0.1"
	shutdownTimeoutSeconds = 15
	httpPort               = "8083"
)

var (
	sigChan     = make(chan os.Signal, 1)
	serverReady = true
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, version)
}

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func main() {
	log.Println("Staring server on port " + httpPort)

	http.HandleFunc("/version", VersionHandler)
	http.HandleFunc("/liveneess", LivenessHandler)

	http.HandleFunc("/readiness/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.URL.Path == "/readiness/ready" {
				serverReady = true
			} else if r.URL.Path == "/readiness/notready" {
				serverReady = false
			}
		}

		if serverReady {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "NOT OK")
		}
	})

	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sigChan <- syscall.SIGTERM
			fmt.Fprintln(w, "Shutdown initiated")
		}
	})

	s := http.Server{
		Addr:              ":" + httpPort,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	//Block until an unterrupt signal is received.
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	log.Printf("Shutting down server in %d seconds...", shutdownTimeoutSeconds)
	serverReady = false
	time.Sleep(shutdownTimeoutSeconds * time.Second)

	err := s.Shutdown(context.Background())
	if err != nil {
		log.Printf("%s\n", err)
	}
	log.Println("BYE")
}
