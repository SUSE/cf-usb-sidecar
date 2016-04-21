package common

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoNewLogger(t *testing.T) {
	testObj := NewLogger("info")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*lager.logger", "Invalid object returned.")
}

func TestDebugNewLogger(t *testing.T) {
	testObj := NewLogger("debug")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*lager.logger", "Invalid object returned.")
}

func TestErrorNewLogger(t *testing.T) {
	testObj := NewLogger("error")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*lager.logger", "Invalid object returned.")
}

func TestFatalNewLogger(t *testing.T) {
	testObj := NewLogger("fatal")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*lager.logger", "Invalid object returned.")
}

func TestDefaultNewLogger(t *testing.T) {
	testObj := NewLogger("anything")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*lager.logger", "Invalid object returned.")
}
