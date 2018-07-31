// Copyright (c) 2018 Vincent Landgraf

package ghostscript

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// NewThumbnailerContext creates a jpeg thumbnail from the given source pdf
func (c Config) NewThumbnailerContext(ctx context.Context, source io.Reader, dpi int, gen WriterGeneratorFunc) error {
	// create a temp file for the pager
	tmpfile, err := ioutil.TempFile("", "ghostscript-thumbnailer.")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// put the source content to the tmp file location
	_, err = io.Copy(tmpfile, source)
	if err != nil {
		return err
	}

	dpiStr := fmt.Sprintf("-r%d", dpi)
	args := []string{"-dPDFFitPage", dpiStr,
		"-dMaxBitmap=500000000", "-dAlignToPixels=0", "-dGridFitTT=2",
		"-sDEVICE=jpeg", "-dTextAlphaBits=4", "-dGraphicsAlphaBits=4",
		"-sDEVICE=pngalpha", tmpfile.Name()}

	return c.NewPagerContext(ctx, args, gen)
}
