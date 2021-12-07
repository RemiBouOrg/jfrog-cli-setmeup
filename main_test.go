package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getApp(t *testing.T) {
	app := getApp()
	assert.Equal(t, 4, len(app.Commands))
}
