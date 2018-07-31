// Copyright (c) 2018 Vincent Landgraf

package ghostscript

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextExtractor(t *testing.T) {
	source, err := os.Open("test/input.pdf")
	assert.NoError(t, err)
	defer source.Close()

	content, err := DefaultConfig.NewTextExtractorContext(context.Background(), source)

	assert.Len(t, content, 2)
	assert.Contains(t, content[0], "www.enbw.com/erneuerbare")
	assert.Contains(t, content[1], "Auswirkungen unserer Stromerzeugung")
}
