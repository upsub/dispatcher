package message

import "testing"

func TestCreate(t *testing.T) {
	msg := Create("testing")

	if msg.Payload != "testing" {
		t.Errorf("message wasn't created correctly")
	}
}

func TestText(t *testing.T) {
	msg := Text("channel", "payload")

	if msg.Header.Get("upsub-channel") != "channel" {
		t.Error("upsub-channel wasn't set correctly")
	}

	if msg.Header.Get("upsub-message-type") != "text" {
		t.Error("upsub-message-type wasn't set correctly")
	}

	if msg.Payload != "payload" {
		t.Error("payload wasn't set correctly")
	}
}

func TestPing(t *testing.T) {
	msg := Ping()

	if msg.Header.Get("upsub-message-type") != "ping" {
		t.Error("ping message wasn't created correctly")
	}
}

func TestPong(t *testing.T) {
	msg := Pong()

	if msg.Header.Get("upsub-message-type") != "pong" {
		t.Error("pong message wasn't created correctly")
	}
}

func TestDecode(t *testing.T) {
	msgString := "{\"headers\":{\"upsub-message-type\":\"text\",\"upsub-channel\":\"channel\"},\"payload\":\"\\\"hello world!\\\"\"}"

	msg, _ := Decode([]byte(msgString))

	if msg.Header.Get("upsub-message-type") != "text" {
		t.Error("upsub message type isn't of type text")
	}

	if msg.Header.Get("upsub-channel") != "channel" {
		t.Error("upsub channel isn't set correctly")
	}

	if msg.Payload != "\"hello world!\"" {
		t.Error("the payload isn't correct")
	}

}

func TestParseBatch(t *testing.T) {
	msgString := "{\"headers\":{\"upsub-message-type\":\"batch\"},\"payload\":\"[{\\\"headers\\\":{\\\"upsub-message-type\\\":\\\"text\\\",\\\"upsub-channel\\\":\\\"channel\\\"},\\\"payload\\\":\\\"payload\\\"},{\\\"headers\\\":{\\\"upsub-message-type\\\":\\\"subscribe\\\"},\\\"payload\\\":\\\"channel\\\"}]\"}"

	msg, _ := Decode([]byte(msgString))

	payload := msg.ParseBatch()

	if payload[0].Header.Get("upsub-message-type") != "text" {
		t.Error("first message wasn't type of text")
	}

	if payload[1].Header.Get("upsub-message-type") != "subscribe" {
		t.Error("second message wasn't type of subscribe")
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
	encoded, _ := Encode(msg)
	if len(encoded) == 0 {
		t.Error("static encode failed")
	}
}

func TestEncode(t *testing.T) {
	msg := Create("payload")
	encoded, _ := msg.Encode()
	if len(encoded) == 0 {
		t.Error("static encode failed")
	}
}
