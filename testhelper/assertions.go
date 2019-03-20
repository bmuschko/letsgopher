package testhelper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func FileNotExists(t *testing.T, path string, msgAndArgs ...interface{}) bool {
	t.Helper()
	info, err := os.Lstat(path)
	if info != nil || err == nil {
		if info.IsDir() {
			return assert.Fail(t, fmt.Sprintf("%q is a directory", path), msgAndArgs...)
		}
		return assert.Fail(t, fmt.Sprintf("found file %q", path), msgAndArgs...)
	}
	return true
}
