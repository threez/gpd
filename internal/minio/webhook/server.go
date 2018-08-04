package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/threez/gpd/pkg/minio/webhook"
)

const basePath = "/gpd/events"

// Server to auto thumbnail and create metadata for pdf files in minio
// buckets
type Server struct {
	Token   string
	Handler webhook.Handler
}

// ListenAndServeTLS starts listening for bucket events like object creation,
// deletion etc. and creates and or removes the corresponding meta data
// and thumbnails. Returns nil if the passed certificate and key are empty.
func (s Server) ListenAndServeTLS(secureAddress, certPath, keyPath string) *http.Server {
	if certPath != "" && keyPath != "" {
		log.Printf("Listen for secure webhooks on: https://%s%s", secureAddress, s.path())
		secureServer := http.Server{
			Addr:    secureAddress,
			Handler: s,
		}
		go func() {
			err := secureServer.ListenAndServeTLS(certPath, keyPath)
			if err != nil && err.Error() != "http: Server closed" {
				log.Fatal(err)
			}
		}()
		return &secureServer
	}
	return nil
}

// ListenAndServe starts listening for bucket events like object creation,
// deletion etc. and creates and or removes the corresponding meta data
// and thumbnails.
func (s Server) ListenAndServe(address string) *http.Server {
	log.Printf("Listen for webhooks on: http://%s%s", address, s.path())
	server := http.Server{
		Addr:    address,
		Handler: s,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err.Error() != "http: Server closed" {
			log.Fatal(err)
		}
	}()
	return &server
}

// ServeHTTP processes incoming webhook requests
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	if r.URL.Path != basePath {
		w.WriteHeader(404)
		return
	}

	// check token
	if r.URL.Query().Get("token") != s.Token {
		w.WriteHeader(401)
		return
	}

	// parse event
	var event webhook.Event
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&event)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Failed to parse event: %s", err)
		return
	}
	log.Printf("WebHook [%s] %s\n", event.EventName, event.Key)

	// process event
	err = s.Handler.ProcessEvent(r.Context(), &event)
	if err != nil {
		w.WriteHeader(501)
		fmt.Fprintf(w, "Failed to process event: %s", err)
		return
	}

	w.WriteHeader(200)
}

func (s Server) path() string {
	return fmt.Sprintf("%s?token=%s", basePath, s.Token)
}
