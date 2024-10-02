package tests

import (
	"os"
	"testing"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestFromConfig(t *testing.T) {
	cfg, err := config.New()
	assert.Nil(t, err)
	assert.Equal(t, cfg.Http.Addr, "0.0.0.0:3000")

	os.Setenv("HTTP_ADDR", ":3000")

	cfg, err = config.New()
	assert.Nil(t, err)
	assert.Equal(t, cfg.Http.Addr, ":3000")
}
