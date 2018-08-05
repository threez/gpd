package ghostscript

import "bytes"

// channelBuffer implements io.WriteCloser, Writes go
// to the buffer and on close the buffer will be sent to
// the channel
type channelBuffer struct {
	Channel chan<- []byte
	bytes.Buffer
}

func (b *channelBuffer) Close() error {
	b.Channel <- b.Bytes()
	return nil
}
