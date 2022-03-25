package headers

type MapHeader struct {
	headers map[string]string
}

func NewMapHeader() *MapHeader {
	headers := make(map[string]string)
	return &MapHeader{
		headers: headers,
	}
}

func (h *MapHeader) GetMap() map[string]string {
	return h.headers
}

func (h *MapHeader) ContentTypeJSON() *MapHeader {
	h.headers["Content-Type"] = "application/json"
	return h
}
