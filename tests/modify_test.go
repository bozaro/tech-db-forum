package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvertCase(t *testing.T) {
	assert.Equal(t, "iwDUVBH693Qw6PyJ", InvertCase("IWduvbh693qW6pYj"))
}
