package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	minio "github.com/minio/minio-go"
	"github.com/threez/gpd/internal/ghostscript"
	"github.com/threez/gpd/internal/minio/webhook"
	gs "github.com/threez/gpd/pkg/ghostscript"
)

var endpoint, accessKeyID, secretAccessKey string
var cert, key, token string
var address = ":3080"
var secureAddress = ":3443"
var dpi = 70
var useSSL = true

func main() {
	parseFlags()

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatal(err)
	}

	ep := ghostscript.EventProcessor{
		Client: client,
		DPI:    dpi,
		GS:     gs.DefaultConfig,
	}

	s := webhook.Server{
		Token:   token,
		Handler: &ep,
	}

	http := s.ListenAndServe(address)
	defer http.Close()

	https := s.ListenAndServeTLS(secureAddress, cert, key)
	defer https.Close()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals // wait until program stop
}
