package message

// Header is a map of the message headers
type Header map[string]string

// Get returns a header keys value
func (h Header) Get(header string) string {
	return h[header]
}

// Set appends or modifies the headers keys value
func (h Header) Set(header string, value string) {
	h[header] = value
}
