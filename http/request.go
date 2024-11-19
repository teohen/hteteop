package http

import (
	"errors"
	"strings"
)

type Request struct {
	Method  string
	URI     string
	Headers map[string]string
}

const CRLF = "\r\n"

func getHeaderEnd(header string) int {

	idx := strings.Index(header, CRLF+CRLF)

	if idx > 0 {
		return idx
	}

	return len(header) - 1
}

func extractURI(uri string) string {
	paramsIdx := strings.Index(uri, "?")

	if paramsIdx < 0 {
		paramsIdx = len(uri)
	}

	return uri[:paramsIdx]
}

func parseRequestLine(requestLine string) (err error, method, path, ptcl string) {
	rl := strings.Split(requestLine, " ")

	uri := extractURI(rl[1])

	if len(rl) == 3 {
		return nil, rl[0], uri, rl[2]
	}
	return errors.New("malformed request line"), "", "", ""
}

func getHeaderValue(header string) string {
	return strings.TrimSpace(strings.ReplaceAll(header, CRLF, ""))
}

func parseHeader(h string) map[string]string {
	headers := make(map[string]string)

	for _, header := range strings.Split(h, "\r\n") {
		splitIdx := strings.Index(header, ":")

		if splitIdx >= 0 {
			key := header[:splitIdx]
			value := getHeaderValue(header[splitIdx+1:])
			headers[key] = value
		}

	}

	return headers
}

func ParseRequest(r string) *Request {
	reqLineIdx := strings.Index(r, CRLF)
	headerEndIdx := getHeaderEnd(r)

	err, method, path, ptcl := parseRequestLine(r[:reqLineIdx])

	if err != nil || ptcl != "HTTP/1.1" {
		return nil
	}

	header := parseHeader(r[reqLineIdx+2 : headerEndIdx])

	http := Request{
		Method:  method,
		URI:     path,
		Headers: header,
	}

	return &http
}
