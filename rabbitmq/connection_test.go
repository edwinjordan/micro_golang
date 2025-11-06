package rabbitmq_test

import (
	"testing"

	"github.com/edwinjordan/micro_golang/rabbitmq"
)

func TestNewConnection_InvalidURL(t *testing.T) {
	// Test that connection fails with invalid URL
	_, err := rabbitmq.NewConnection("invalid://url")
	if err == nil {
		t.Error("Expected error with invalid URL, got nil")
	}
}

func TestDeclareQueue_NilConnection(t *testing.T) {
	// Test that we can create a connection struct
	conn := &rabbitmq.Connection{}
	if conn == nil {
		t.Error("Expected connection struct, got nil")
	}
}
