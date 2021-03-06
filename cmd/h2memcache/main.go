package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/coocood/freecache"
	"golang.org/x/crypto/acme/autocert"
)

var cache *freecache.Cache
var apikey string
var certManager *autocert.Manager

func main() {
	cacheSize, _ := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	hostname := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")
	cert := os.Getenv("TLS_CERT_FILE")
	key := os.Getenv("TLS_KEY_FILE")
	certDir := os.Getenv("TLS_CERT_DIR")
	apikey = os.Getenv("API_KEY")

	certManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostname),
		Cache:      autocert.DirCache(certDir),
	}

	cache = freecache.NewCache(cacheSize * 1024 * 1024)

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	timeout := 5 * time.Second

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: &handler{},
	}

	ssl := false

	if cert != "" && key != "" {
		ssl = true
	} else if hostname != "" {
		srv.TLSConfig = certManager.TLSConfig()
		ssl = true
	}

	go func() {
		if ssl == false {
			log.Printf("Listening on port %s", port)
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		} else {
			log.Printf("Listening on port %s (TLS)", port)
			if err := srv.ListenAndServeTLS(cert, key); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("Shutdown with timeout: %s\n", timeout)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}
