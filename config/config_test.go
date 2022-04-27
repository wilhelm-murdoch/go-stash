package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wilhelm-murdoch/go-stash/config"
)

func TestConfigNew(t *testing.T) {
	_, err := config.New("../.stash.yaml")
	assert.Nil(t, err, "Expected no errors, but got %s", err)
}
