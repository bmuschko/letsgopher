package environment

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSettings(t *testing.T) {
	settings := Settings
	assert.True(t, strings.HasSuffix(settings.Home.String(), ".letsgopher"))
}
