package tests

import (
	"os"
	"testing"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_From_Config(t *testing.T) {
	cfg, err := config.New(logger.New())
	assert.Nil(t, err)
	assert.Equal(t, cfg.Http.Addr, "0.0.0.0:3000")

	os.Setenv("DPOINT_HTTP_ADDR", ":3000")

	cfg, err = config.New(logger.New())
	assert.Nil(t, err)
	assert.Equal(t, cfg.Http.Addr, ":3000")
}
