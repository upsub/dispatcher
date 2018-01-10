package message

// Header is a map of the message headers
type Header map[string]string

// Get returns a header
func (h Header) Get(header string) string {
	value := h[header]

	return value
}

// Set appends or modifies the headers
func (h Header) Set(header string, value string) {
	h[header] = value
}
