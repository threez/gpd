// Copyright (c) 2018 Vincent Landgraf

package ghostscript

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
)

var continueExp = []byte(">>showpage, press <return> to continue<<\n")
var continueExpLen = len(continueExp)

// WriterGeneratorFunc function that generates a writer for each page
type WriterGeneratorFunc func(ctx context.Context, page int) (w io.WriteCloser, err error)

// NewPagerContext uses the given builder to generate writers for each page
// Note: pages start at 1 not 0
func (c Config) NewPagerContext(ctx context.Context, args []string, generateWriter WriterGeneratorFunc) error {
	cmd := exec.CommandContext(ctx, c.Command, append(c.Arguments, args...)...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = stderr.Close()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	sob := bufio.NewReader(stdout)
	i := 0
	var w io.WriteCloser

	for {
		buf, err := sob.ReadSlice('\n')
		n := len(buf)

		if err == io.EOF && n == 0 {
			break
		}

		// create new writer using generator function
		if w == nil {
			i++
			w, err = generateWriter(ctx, i)
			if err != nil {
				return err
			}
		}

		// Check if line is continue expression
		if n >= continueExpLen && bytes.Compare(buf[n-continueExpLen:], continueExp) == 0 {
			w.Write(buf[:n-continueExpLen])

			err = w.Close()
			if err != nil {
				return err
			}

			// page to next page
			_, err = stdin.Write([]byte{0x0A})
			if err != nil {
				return err
			}

			w = nil
		} else {
			// otherwise write the data into the current writer
			w.Write(buf)
		}

		// check on read errors
		if err == io.EOF {
			break
		}
	}

	if w != nil {
		err = w.Close()
		if err != nil {
			return err
		}
	}

	err = stdin.Close()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
