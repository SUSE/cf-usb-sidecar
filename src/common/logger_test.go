package common

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoNewLogger(t *testing.T) {
	testObj := NewLogger("info", "test-logger")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*logrus.Logger", "Invalid object returned.")
}

func TestDebugNewLogger(t *testing.T) {
	testObj := NewLogger("debug", "test-logger")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*logrus.Logger", "Invalid object returned.")
}

func TestErrorNewLogger(t *testing.T) {
	testObj := NewLogger("error", "test-logger")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*logrus.Logger", "Invalid object returned.")
}

func TestFatalNewLogger(t *testing.T) {
	testObj := NewLogger("fatal", "test-logger")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*logrus.Logger", "Invalid object returned.")
}

func TestDefaultNewLogger(t *testing.T) {
	testObj := NewLogger("anything", "test-logger")

	assert.NotNil(t, testObj, "Unexpected error.")
	assert.Equal(t, reflect.TypeOf(testObj).String(), "*logrus.Logger", "Invalid object returned.")
}
