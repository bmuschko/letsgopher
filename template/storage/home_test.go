package storage

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestHome(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	home := Home(tmpHome)
	assert.Equal(t, tmpHome, home.String())
	assert.Equal(t, filepath.Join(tmpHome, "archive"), home.ArchiveDir())
	assert.Equal(t, filepath.Join(tmpHome, "templates.yaml"), home.TemplatesFile())
}
