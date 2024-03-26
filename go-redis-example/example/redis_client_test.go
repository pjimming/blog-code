package example

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	ast := assert.New(t)
	ast.NotNil(RedisCli)
	ast.Nil(RedisCli.Ping(context.Background()).Err())
}
