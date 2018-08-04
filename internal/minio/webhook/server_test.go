package webhook

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threez/gpd/pkg/minio/webhook"
)

var testServer = Server{
	Token: "test-token",
}

func TestServerPath(t *testing.T) {
	s := testServer.ListenAndServe("localhost:9090")

	r, err := http.Get(fmt.Sprintf("http://%s/some-url", s.Addr))
	assert.NoError(t, err)
	assert.Equal(t, 404, r.StatusCode)
	assert.NoError(t, s.Close())
}

func TestTokenPath(t *testing.T) {
	s := testServer.ListenAndServe("localhost:9091")

	r, err := http.Get(fmt.Sprintf("http://%s%s", s.Addr, basePath))
	assert.NoError(t, err)
	assert.Equal(t, 401, r.StatusCode)
	assert.NoError(t, s.Close())
}

func TestParsingPath(t *testing.T) {
	s := testServer.ListenAndServe("localhost:9092")

	r, err := http.Post(fmt.Sprintf("http://%s%s?token=%s", s.Addr, basePath, testServer.Token),
		"application/json", bytes.NewBufferString(""))
	assert.NoError(t, err)
	assert.Equal(t, 400, r.StatusCode)
	assert.NoError(t, s.Close())
}

type errHandler struct{}

func (e errHandler) ProcessEvent(ctx context.Context, ev *webhook.Event) error {
	return fmt.Errorf("Failed")
}

func TestProcessingPath(t *testing.T) {
	testServer.Handler = &errHandler{}
	s := testServer.ListenAndServe("localhost:9093")

	r, err := http.Post(fmt.Sprintf("http://%s%s?token=%s", s.Addr, basePath, testServer.Token),
		"application/json", bytes.NewBufferString(`{"EventName": "s3:ObjectCreated:Put", "Key": "path/or/key/some-document.pdf"}`))
	assert.NoError(t, err)
	assert.Equal(t, 501, r.StatusCode)
	assert.NoError(t, s.Close())
}

type simpleRecorder struct{ ev *webhook.Event }

func (e *simpleRecorder) ProcessEvent(ctx context.Context, ev *webhook.Event) error {
	e.ev = ev
	return nil
}

func TestWebhookHappyPath(t *testing.T) {
	rec := &simpleRecorder{}
	testServer.Handler = rec
	s := testServer.ListenAndServe("localhost:9094")

	r, err := http.Post(fmt.Sprintf("http://%s%s?token=%s", s.Addr, basePath, testServer.Token),
		"application/json", bytes.NewBufferString(`{"EventName": "s3:ObjectCreated:Put", "Key": "path/or/key/some-document.pdf"}`))
	assert.NoError(t, err)
	assert.Equal(t, 200, r.StatusCode)
	assert.NotNil(t, rec.ev)
	assert.Equal(t, "s3:ObjectCreated:Put", rec.ev.EventName)
	assert.Equal(t, "path/or/key/some-document.pdf", rec.ev.Key)
	assert.NoError(t, s.Close())
}
