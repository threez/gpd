// Copyright (c) 2018 Vincent Landgraf

package ghostscript

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThumbnailer(t *testing.T) {
	gen := func(ctx context.Context, page int) (io.WriteCloser, error) {
		file := filepath.Join("test", fmt.Sprintf("gen-%d.jpg", page))
		return os.Create(file)
	}

	source, err := os.Open("test/input.pdf")
	assert.NoError(t, err)
	defer source.Close()

	err = DefaultConfig.NewThumbnailerContext(context.Background(), source, 70, gen)
	assert.NoError(t, err)

	assert.FileExists(t, "test/gen-1.jpg")
	md5, err := md5File("test/gen-1.jpg")
	assert.NoError(t, err)
	assert.Equal(t, "76e983bf9a247a84f706d74695de03f6", md5)

	assert.FileExists(t, "test/gen-2.jpg")
	md5, err = md5File("test/gen-2.jpg")
	assert.NoError(t, err)
	assert.Equal(t, "45dc970bd4764b8295cbc71816deb585", md5)

	os.Remove("test/gen-1.jpg")
	os.Remove("test/gen-2.jpg")
}

func md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
