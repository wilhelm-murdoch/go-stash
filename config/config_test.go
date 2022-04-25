package config_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wilhelm-murdoch/go-stash/config"
)

func TestConfigNew(t *testing.T) {
	c, err := config.New(".stash.test.yaml")
	assert.Nil(t, err, "Expected no errors, but got %s", err)
	fmt.Printf("c: %v\n", c)
}
