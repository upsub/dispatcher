package message

import (
	"testing"
)

func TestCreate(t *testing.T) {
	msg := Create("testing")

	if msg.Payload != "testing" {
		t.Errorf("message wasn't created correctly")
	}
}

func TestText(t *testing.T) {
	msg := Text("channel", "payload")

	if msg.Channel != "channel" {
		t.Error("upsub-channel wasn't set correctly")
	}

	if msg.Type != "text" {
		t.Error("upsub-message-type wasn't set correctly")
	}

	if msg.Payload != "payload" {
		t.Error("payload wasn't set correctly")
	}
}

func TestPing(t *testing.T) {
	msg := Ping()

	if msg.Type != "ping" {
		t.Error("ping message wasn't created correctly")
	}
}

func TestPong(t *testing.T) {
	msg := Pong()

	if msg.Type != "pong" {
		t.Error("pong message wasn't created correctly")
	}
}

func TestResponseAction(t *testing.T) {
	msg := ResponseAction([]string{"channel", "channel-2/event"}, "action")

	if msg.Type != "text" {
		t.Error("SubscribeResponse upsub-message-type wasn't correct type")
	}

	if msg.Channel != "channel:action,channel-2/event:action" {
		t.Error("SubscribeResponse upsub-channel wasn't set correctly")
	}

	if msg.Payload != "channel:action,channel-2/event:action" {
		t.Error("SubscribeResponse upsub-channel wasn't set correctly")
	}
}

func TestDecode(t *testing.T) {
	msgString := "text channel\n\nhello world!"

	msg, _ := Decode([]byte(msgString))

	if msg.Type != "text" {
		t.Error("upsub message type isn't of type text")
	}

	if msg.Channel != "channel" {
		t.Error("upsub channel isn't set correctly")
	}

	if msg.Payload != "hello world!" {
		t.Error("the payload isn't correct")
	}

}

func TestParseBatch(t *testing.T) {
	msgString := "batch\n\n[\"text channel\\n\\npayload\", \"subscribe\\n\\nchannel\"]"

	msg, _ := Decode([]byte(msgString))

	payload := msg.ParseBatch()

	if payload[0].Type != "text" {
		t.Error("first message wasn't type of text")
	}

	if payload[0].Payload != "payload" {
		t.Error("Payload wasn't parsed correctly")
	}

	if payload[1].Type != "subscribe" {
		t.Error("second message wasn't type of subscribe")
	}

	if payload[1].Payload != "channel" {
		t.Error("Payload wasn't parsed correctly")
	}

}

func TestParseBatchFail(t *testing.T) {
	msgString := "{\"headers\":{\"upsub-message-type\":\"batch\"},\"payload\":\"[{\\\"headers:{\\\"upsub-message-type\\\":\\\"text\\\",\\\"upsub-channel\\\":\\\"channel\\\"},\\\"payload\\\":\\\"payload\\\"},{\\\"headers\\\":{\\\"upsub-message-type\\\":\\\"subscribe\\\"},\\\"payload\\\":\\\"channel\\\"}]\"}"

	msg, _ := Decode([]byte(msgString))

	payload := msg.ParseBatch()

	if len(payload) > 0 {
		t.Error("didn't handle parsing of batch messages correctly")
	}
}

func TestStaticEncode(t *testing.T) {
	msg := Create("payload")
	encoded := Encode(msg)
	if len(encoded) == 0 {
		t.Error("static encode failed")
	}
}

func TestEncode(t *testing.T) {
	msg := Create("payload")
	encoded := msg.Encode()
	if len(encoded) == 0 {
		t.Error("static encode failed")
	}
}
