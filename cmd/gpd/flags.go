package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseFlags() {
	flag.StringVar(&token, "token", os.Getenv("GPD_TOKEN"), `Token to securely verify the validity of webhooks
		Environment variable: GPD_TOKEN`)

	// DPI
	if newDPI, ok := os.LookupEnv("GPD_DPI"); ok {
		var err error
		dpi, err = strconv.Atoi(newDPI)
		if err != nil {
			log.Fatalf("Invalid DPI: %s", err)
		}
	}
	flag.IntVar(&dpi, "dpi", dpi, `DPI resolution for thumbnail versions
		Environment variable: GPD_DPI`)

	// HTTP
	if newAddress, ok := os.LookupEnv("GPD_ADDRESS"); ok {
		address = newAddress
	}
	flag.StringVar(&address, "address", address, `Server port for minio webhooks service
		Environment variable: GPD_ADDRESS`)

	// HTTPS
	if newSecureAddress, ok := os.LookupEnv("GPD_SECURE_ADDRESS"); ok {
		secureAddress = newSecureAddress
	}
	flag.StringVar(&secureAddress, "secure-address", secureAddress, `Server port for the secure minio webhooks service
		Environment variable: GPD_SECURE_ADDRESS`)
	flag.StringVar(&cert, "cert", os.Getenv("GPD_CERT"), `Path to the certificate file (PEM) for the secure webhook service
		Environment variable: GPD_CERT`)
	flag.StringVar(&key, "key", os.Getenv("GPD_KEY"), `Path to the provate key file (PEM) for the secure webhook service
		Environment variable: GPD_KEY`)

	// MINIO
	flag.StringVar(&endpoint, "endpoint", os.Getenv("GPD_ENDPOINT"), `Endpoint of the cloud provider (host + port allowed)
		Environment variable: GPD_ENDPOINT`)
	flag.StringVar(&accessKeyID, "access-key-id", os.Getenv("GPD_ACCESS_KEY_ID"), `Minio credentials (ID)
		Environment variable: GPD_ACCESS_KEY_ID`)
	flag.StringVar(&secretAccessKey, "secret-access-key", os.Getenv("GPD_SECRET_ACCESS_KEY"), `Minio credentials (secret)
		Environment variable: GPD_SECRET_ACCESS_KEY`)

	noSSL := flag.Bool("no-ssl", strings.ToLower(os.Getenv("GPD_NO_SSL")) == "yes",
		`disables ssl on minio connections
	Environment variable: GPD_NO_SSL (YES|NO)`)

	flag.Parse()

	useSSL = !(noSSL != nil && *noSSL == true)
}
