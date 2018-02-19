package message

import "testing"

func TestMessageTypes(t *testing.T) {
	if TEXT != "text" {
		t.Error("TEXT type not correct")
	}

	if BATCH != "batch" {
		t.Error("BATCH type not correct")
	}

	if SUBSCRIBE != "subscribe" {
		t.Error("SUBSCRIBE type not correct")
	}

	if UNSUBSCRIBE != "unsubscribe" {
		t.Error("UNSUBSCRIBE type not correct")
	}

	if PING != "ping" {
		t.Error("PING type not correct")
	}

	if PONG != "pong" {
		t.Error("PONG type not correct")
	}
}
