package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	ast := assert.New(t)
	ast.NotNil(RedisCli)
}
