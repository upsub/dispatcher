package message

import "testing"

func TestGet(t *testing.T) {
	header := &Header{"key": "value"}

	if header.Get("key") != "value" {
		t.Error("Failed at retrieving value from key")
	}
}

func TestSet(t *testing.T) {
	header := &Header{}
	header.Set("key", "value")

	if header.Get("key") != "value" {
		t.Error("Failed at setting a new header key with a value")
	}
}
