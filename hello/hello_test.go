package hello

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGiveUCaseOutput(t *testing.T) {
	text := uCaseOutput("hello World")
	assert.Equal(t, "HELLO WORLD", text)
}

func TestGiveUCaseOutputNeg(t *testing.T) {
	text := uCaseOutput("hello World")
	assert.NotEqual(t, "hello world", text)
}
