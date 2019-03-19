package cmd

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUndefinedVersion(t *testing.T) {
	b := bytes.NewBuffer(nil)
	version := &versionCmd{
		out: b,
	}
	err := version.run()

	assert.Nil(t, err)
	assert.Equal(t, "letsgopher \n", b.String())
}

func TestSemanticVersion(t *testing.T) {
	b := bytes.NewBuffer(nil)
	SetVersion("1.2.3")
	version := &versionCmd{
		out: b,
	}
	err := version.run()

	assert.Nil(t, err)
	assert.Equal(t, "letsgopher 1.2.3\n", b.String())
}

func TestVersionForError(t *testing.T) {
	b := new(WriterMock)
	version := &versionCmd{
		out: b,
	}
	err := version.run()

	assert.NotNil(t, err)
	assert.Equal(t, "expected", err.Error())
}

type WriterMock struct {
	mock.Mock
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	return 0, errors.New("expected")
}
