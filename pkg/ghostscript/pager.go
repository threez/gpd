// Copyright (c) 2018 Vincent Landgraf

package ghostscript

import (
	"bufio"
	"context"
	"io"
	"os/exec"
)

const continueExp = ">>showpage, press <return> to continue<<"

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
	var buf [4096]byte
	i := 0
	var w io.WriteCloser

	for {
		n, err := sob.Read(buf[:])
		if n > 0 && n > continueExpLen {
			if w == nil {
				i++
				w, err = generateWriter(ctx, i)
				if err != nil {
					return err
				}
			}

			if string(buf[n-continueExpLen-1:n-1]) == continueExp {
				w.Write(buf[:n-continueExpLen-1])
				err := w.Close()
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
				w.Write(buf[:n])
			}
		}

		if err == io.EOF {
			break
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
