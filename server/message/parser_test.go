package message

import (
	"testing"
)

func TestParse(t *testing.T) {
	str := "text channel\n\nhello"
	msg, _ := parse([]byte(str))

	if msg.Type != TEXT {
		t.Error("Message.Type wasn't parsed correctly")
	}

	if msg.Channel != "channel" {
		t.Error("Message.Channel wasn't parsed correctly")
	}

	if msg.Payload != "hello" {
		t.Error("Message.Channel wasn't parsed correctly")
	}
}

func TestParseShouldRemoveLeadingSpaces(t *testing.T) {
	str := "text channel\n\n    hello"
	msg, _ := parse([]byte(str))

	if msg.Type != TEXT {
		t.Error("Message.Type wasn't parsed correctly")
	}

	if msg.Channel != "channel" {
		t.Error("Message.Channel wasn't parsed correctly")
	}

	if msg.Payload != "hello" {
		t.Error("Message.Channel wasn't parsed correctly")
	}
}

func TestParseShouldCreateHeaderMap(t *testing.T) {
	str := `text channel
header-key: header-value

hello`
	msg, _ := parse([]byte(str))

	if msg.Type != TEXT {
		t.Error("Message.Type wasn't parsed correctly")
	}

	if msg.Channel != "channel" {
		t.Error("Message.Channel wasn't parsed correctly")
	}

	if msg.Header.Get("header-key") != "header-value" {
		t.Error("Message.Header wasn't parsed correctly")
	}

	if msg.Payload != "hello" {
		t.Error("Message.Payload wasn't parsed correctly")
	}
}
