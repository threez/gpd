package ghostscript

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"

	minio "github.com/minio/minio-go"
	"github.com/threez/gpd/pkg/ghostscript"
	"github.com/threez/gpd/pkg/minio/webhook"
)

// EventProcessor processes the incoming webhook events from minio
// and uses the ghostscript package to create text and thumbnail versions
// of the PDF documents
type EventProcessor struct {
	Client *minio.Client
	DPI    int
	GS     ghostscript.Config
}

// ProcessEvent entrypoint for the webhooks
func (p *EventProcessor) ProcessEvent(ctx context.Context, e *webhook.Event) error {
	for _, record := range e.Records {
		if strings.HasPrefix(e.EventName, "s3:ObjectCreated:") {
			err := p.processCreateEvent(ctx, record.S3.Bucket.Name, record.S3.Object.Key)
			if err != nil {
				return fmt.Errorf("Failed to process NEW %s/%s: %s", record.S3.Bucket.Name, record.S3.Object.Key, err)
			}
		} else if strings.HasPrefix(e.EventName, "s3:ObjectRemoved:") {
			err := p.processRemoveEvent(ctx, record.S3.Bucket.Name, record.S3.Object.Key)
			if err != nil {
				return fmt.Errorf("Failed to process DEL %s/%s: %s", record.S3.Bucket.Name, record.S3.Object.Key, err)
			}
		}
	}

	return nil
}

func (p *EventProcessor) processCreateEvent(ctx context.Context, bucketName, key string) error {
	document, err := p.fetchDocument(ctx, bucketName, key)
	if err != nil {
		return err
	}

	log.Printf("createThumbnails %s/%s", bucketName, key)
	err = p.createThumbnails(context.Background(), document, bucketName, key)
	if err != nil {
		return err
	}

	log.Printf("createJSONDocument %s/%s", bucketName, key)
	err = p.createJSONDocument(ctx, document, bucketName, key)
	if err != nil {
		return err
	}
	log.Printf("Done %s/%s", bucketName, key)

	return nil
}

func (p *EventProcessor) processRemoveEvent(ctx context.Context, bucketName, key string) error {
	doneCh := make(chan struct{})
	defer close(doneCh) // stop list objects
	infos := p.Client.ListObjectsV2(bucketName, keyPathNoSuffix(key), true, doneCh)

	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				return err
			}
			return nil
		case info, ok := <-infos:
			if !ok { // no more infos to process
				return nil
			}

			// Search for gpd artifacts
			if strings.Contains(filepath.Base(info.Key), ".gpd") {
				log.Printf("Delete gpd artifact: %s", info.Key)
				err := p.Client.RemoveObject(bucketName, info.Key)
				if err != nil {
					return err
				}
			}
		}
	}
}

func (p *EventProcessor) fetchDocument(ctx context.Context, bucketName, key string) ([]byte, error) {
	// Request Object from S3
	object, err := p.Client.GetObjectWithContext(ctx, bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetObject failed: %s", err)
	}

	// Cache it in memory
	document, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("Read failed: %s", err)
	}

	return document, nil
}

func (p *EventProcessor) createThumbnails(ctx context.Context, document []byte, bucketName, key string) error {
	// This section implements a buffer producer and consumer.
	// Using streams is not working, as the minio/S3 API want's to
	// know the size of the upload upfront.
	// Therefore a channel of thumbnail buffers will be created and used
	// to communicate between the thumbnail generation and uploading
	// of the same
	chBuffers := make(chan []byte, 1)
	defer close(chBuffers) // close upload go routine afterwards

	// wait until all thumbnails are created
	var wg sync.WaitGroup

	// Generate thumbnails
	gen := func(ctx context.Context, page int) (io.WriteCloser, error) {
		wg.Add(1)
		return &channelBuffer{Channel: chBuffers}, nil
	}

	// Upload thumbnails
	go func() {
		pageNr := 1
		for thumbnailBuffer := range chBuffers {
			r := bytes.NewReader(thumbnailBuffer)
			_, err := p.Client.PutObjectWithContext(ctx, bucketName, pagePath(key, pageNr),
				r, int64(len(thumbnailBuffer)),
				minio.PutObjectOptions{
					ContentType:        "image/png",
					ContentDisposition: "attachment",
				})
			wg.Done()
			pageNr++
			if err != nil {
				log.Printf("Failed to store thumbnail %s/%s", bucketName, pagePath(key, pageNr))
			}
		}
	}()

	err := p.GS.NewThumbnailerContext(ctx, bytes.NewReader(document), p.DPI, gen)
	wg.Wait()
	return err
}

func (p *EventProcessor) createJSONDocument(ctx context.Context, document []byte, bucketName, key string) error {
	buf := bytes.NewBuffer(document)
	// extract pages from document
	pages, err := p.GS.NewTextExtractorContext(ctx, buf)
	if err != nil {
		return err
	}

	// build json data structure
	doc := NewDocument(len(pages))
	for i, pageText := range pages {
		doc.Pages[i].AddText(pageText, 0, 0, 1, 1)
		doc.Pages[i].ThumbnailKey = pagePath(key, i+1)
	}

	// serialize data
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}

	// upload data
	_, err = p.Client.PutObjectWithContext(ctx, bucketName, jsonPath(key),
		bytes.NewReader(out), int64(len(out)), minio.PutObjectOptions{
			ContentType:        "application/json",
			ContentDisposition: "attachment",
		})
	if err != nil {
		log.Printf("Stored new JSON meta data file %s/%s", bucketName, jsonPath(key))
	}
	return err
}

func pagePath(key string, page int) string {
	return generatePath(key, fmt.Sprintf("-%d.gpd.png", page))
}

func jsonPath(key string) string {
	return generatePath(key, ".gpd.json")
}

func generatePath(key, suffix string) string {
	return keyPathNoSuffix(key) + suffix
}

func keyPathNoSuffix(key string) string {
	return strings.TrimSuffix(key, ".pdf")
}
