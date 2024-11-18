package http

import (
	"errors"
	"strings"
)

type HttpRequest struct {
	Method  string
	URI     string
	Headers map[string]string
	Data    string
}

const CRLF = "\r\n"

func getHeaderEnd(header string) int {

	idx := strings.Index(header, CRLF+CRLF)

	if idx > 0 {
		return idx
	}

	return len(header) - 1
}

func parseRequestLine(requestLine string) (err error, method, path, ptcl string) {
	rl := strings.Split(requestLine, " ")

	if len(rl) == 3 {
		return nil, rl[0], rl[1], rl[2]
	}
	return errors.New("malformed request line"), "", "", ""
}

func parseHeader(h string) map[string]string {
	headers := make(map[string]string)

	for _, header := range strings.Split(h, "\r\n") {
		splitIdx := strings.Index(header, ":")

		if splitIdx >= 0 {
			key := header[:splitIdx]
			value := header[splitIdx+1:]
			headers[key] = value
		}

	}

	return headers
}

func ParseRequest(r string) *HttpRequest {
	reqLineIdx := strings.Index(r, CRLF)
	headerEndIdx := getHeaderEnd(r)

	err, method, path, ptcl := parseRequestLine(r[:reqLineIdx])

	if err != nil || ptcl != "HTTP/1.1" {
		return nil
	}

	h := parseHeader(r[reqLineIdx+2 : headerEndIdx])

	http := HttpRequest{
		Method:  method,
		URI:     path,
		Headers: h,
	}

	return &http
}
