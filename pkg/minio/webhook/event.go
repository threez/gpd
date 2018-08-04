package webhook

import (
	"time"
)

// Event high level minio bucket change event
type Event struct {
	EventName string   `json:"EventName"` // "s3:ObjectCreated:Put"
	Key       string   `json:"Key"`       // "path/or/key/some-document.pdf"
	Records   []Record `json:"Records"`
}

// Record details the bucket change event
type Record struct {
	EventVersion      string            `json:"eventVersion"` // "2.0"
	EventSource       string            `json:"eventSource"`  // "minio:s3"
	AwsRegion         string            `json:"awsRegion"`    // "eu-central-1"
	EventTime         time.Time         `json:"eventTime"`    // "2018-08-03T22:05:11Z"
	EventName         string            `json:"eventName"`    // "s3:ObjectCreated:Put"
	UserIdentity      UserIdentity      `json:"userIdentity"`
	RequestParameters map[string]string `json:"requestParameters"`
	ResponseElements  map[string]string `json:"responseElements"`
	S3                S3                `json:"s3"`
	Source            Source            `json:"source"`
}

// UserIdentity captures the accessKey ID / principalId
type UserIdentity struct {
	PrincipalID string `json:"principalId"`
}

// S3 details regarding bucket and object changed
type S3 struct {
	S3SchemaVersion string `json:"s3SchemaVersion"` // "1.0"
	ConfigurationID string `json:"configurationId"` // "Config"
	Bucket          Bucket `json:"bucket"`
	Object          Object `json:"object"`
}

// Bucket details
type Bucket struct {
	Name          string       `json:"name"` // "test"
	OwnerIdentity UserIdentity `json:"ownerIdentity"`
	ARN           string       `json:"arn"` // "arn:aws:s3:::test"
}

// Object details
type Object struct {
	Key          string            `json:"key"`         // "2018_07_06_20_39_23_OCR.pdf"
	Size         uint64            `json:"size"`        // 185162
	ETag         string            `json:"eTag"`        // "196c5061dcde8a73458b69eaab836b5f"
	ContentType  string            `json:"contentType"` // "196c5061dcde8a73458b69eaab836b5f"
	UserMetadata map[string]string `json:"userMetadata"`
	VersionID    string            `json:"versionId"` // "1"
	Sequencer    string            `json:"sequencer"` // "15477F11753CD2E0"
}

// Source in case of copies
type Source struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	UserAgent string `json:"userAgent"`
}
