package main

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"reflect"
	"testing"
)

func Test_getApp(t *testing.T) {
	tests := []struct {
		name string
		want components.App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getApp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCommands(t *testing.T) {
	tests := []struct {
		name string
		want []components.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCommands(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}
