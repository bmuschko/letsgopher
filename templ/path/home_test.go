package path

import (
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestHome(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	home := Home(tmpHome)
	assert.Equal(t, tmpHome, home.String())
	assert.Equal(t, filepath.Join(tmpHome, "archive"), home.ArchiveDir())
	assert.Equal(t, filepath.Join(tmpHome, "templates.yaml"), home.TemplatesFile())
}
