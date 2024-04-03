package main

import (
	"log"
	"net/http"
	"golang.org/x/net/webdav"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"
)

func main() {
	addr := flag.String("addr", ":8443", "Address to listen on")
	certFile := flag.String("cert", "cert.pem", "Path to SSL certificate file")
	keyFile := flag.String("key", "key.pem", "Path to SSL key file")
	davDir := flag.String("davdir", ".", "Directory to serve")
	user := flag.String("user", "admin", "Username for Basic Authentication")
	pass := flag.String("pass", "admin", "Password for Basic Authentication")
	flag.Parse()

	//  logging
	log.Println("Starting server...")

	//  starting server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != *user || p != *pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		davHandler := &webdav.Handler{
			FileSystem: webdav.Dir(*davDir),
			LockSystem: webdav.NewMemLS(),
		}
		davHandler.ServeHTTP(w, r)
	})

	// error handling and server startup
	server := &http.Server{Addr: *addr, Handler: nil}

	// Handling shutdown more or less gracefully
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		log.Println("Shutting down the server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	// error handling on server start
	if err := server.ListenAndServeTLS(*certFile, *keyFile); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
