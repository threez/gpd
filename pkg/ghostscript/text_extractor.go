// Copyright (c) 2018 Vincent Landgraf

package ghostscript

import (
	"context"
	"io"
	"io/ioutil"
	"os"
)

// NewTextExtractorContext extracts the text from the given source pdf
func (c Config) NewTextExtractorContext(ctx context.Context, source io.Reader) ([]string, error) {
	var pages []string

	// create a temp file for the pager
	tmpfile, err := ioutil.TempFile("", "ghostscript-text-extractor.")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// put the source content to the tmp file location
	_, err = io.Copy(tmpfile, source)
	if err != nil {
		return nil, err
	}

	args := []string{"-sDEVICE=txtwrite", tmpfile.Name()}

	err = c.NewPagerContext(ctx, args, func(ctx context.Context, page int) (io.WriteCloser, error) {
		r, w := io.Pipe()

		go func() {
			data, err := ioutil.ReadAll(r)
			if err != nil {
				return
			}
			pages = append(pages, string(data))
			r.Close()
		}()

		return w, nil
	})
	if err != nil {
		return nil, err
	}

	return pages, nil
}
