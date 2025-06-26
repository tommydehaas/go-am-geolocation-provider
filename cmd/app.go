package main

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/tommydehaas/go-am-geolocation-provider/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{ip}", service.GetCountyCodeByIp)

	server := &http.Server{
		Addr: "0.0.0.0:8443",
		TLSConfig: &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := server.ListenAndServeTLS("tls/cert.pem", "tls/key.pem"); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	<-stop
	server.Shutdown(context.Background())
}
