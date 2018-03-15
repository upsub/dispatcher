package message

import (
	"regexp"
	"strings"
)

func extractHeadAndBody(msg string) (string, string) {
	msgNormalized := regexp.MustCompile("\r\n").ReplaceAllString(msg, "\n")
	sections := regexp.MustCompile(`\n\s*\n`).Split(msgNormalized, -1)

	if len(sections) == 1 {
		return sections[0], ""
	}

	return sections[0], strings.TrimSpace(sections[1])
}

func extractTypeAndChannel(str string) (string, string) {
	parts := strings.Split(str, " ")

	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], parts[1]
}

func parseHeader(str string) (string, string) {
	parts := strings.Split(str, ":")
	return parts[0], strings.TrimSpace(strings.Replace(parts[1], "\n", "", 1))
}

func parse(msg []byte) (*Message, error) {
	head, body := extractHeadAndBody(string(msg))
	msgType := ""
	msgChannel := ""
	msgHeader := Header{}

	for i, value := range strings.Split(head, "\n") {
		if i == 0 {
			msgType, msgChannel = extractTypeAndChannel(value)
			continue
		}

		key, value := parseHeader(value)
		msgHeader.Set(key, value)
	}

	return &Message{
		msgType,
		msgChannel,
		msgHeader,
		body,
		false,
	}, nil
}
