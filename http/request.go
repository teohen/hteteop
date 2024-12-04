package http

import (
	"errors"
	"strings"
)

type Request struct {
	Method      string
	URI         string
	Headers     map[string]string
	QueryParams map[string]string
	PathParams  map[string]string
}

const CRLF = "\r\n"

func sanitize(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, CRLF, ""))
}

func getHeaderEnd(header string) int {

	idx := strings.Index(header, CRLF+CRLF)

	if idx > 0 {
		return idx
	}

	return len(header) - 1
}

func parseURI(uri string) (string, map[string]string) {
	paramsIdx := strings.Index(uri, "?")

	params := make(map[string]string)

	if paramsIdx < 0 {
		return sanitize(uri), nil
	}

	paramsStr := uri[paramsIdx+1:]

	for _, p := range strings.Split(paramsStr, "&") {
		kv := strings.Split(p, "=")
		params[sanitize(kv[0])] = sanitize(kv[1])
	}

	return sanitize(uri[:paramsIdx]), params
}

func parseRequestLine(requestLine string) (err error, method, path, ptcl string, params map[string]string) {
	rl := strings.Split(requestLine, " ")

	method = rl[0]
	uri := rl[1]
	ptcl = rl[2]
	params = make(map[string]string)

	path, params = parseURI(uri)

	if len(rl) == 3 {
		return nil, method, path, ptcl, params
	}

	return errors.New("malformed request line"), "", "", "", nil
}

func getHeaderValue(header string) string {
	return sanitize(header)
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

func (r *Request) ParsePathParams(routeUri []string) {
	r.PathParams = make(map[string]string)
	listURI := strings.Split(r.URI, "/")
	for i, part := range routeUri {
		if part != "" && part[0] == '{' {
			key := part[1 : len(part)-1]
			r.PathParams[key] = listURI[i]
		}
	}

}

func ParseRequest(r string) *Request {
	reqLineIdx := strings.Index(r, CRLF)
	headerEndIdx := getHeaderEnd(r)

	err, method, path, ptcl, params := parseRequestLine(r[:reqLineIdx])

	if err != nil || ptcl != "HTTP/1.1" {
		return nil
	}

	header := parseHeader(r[reqLineIdx+2 : headerEndIdx])

	request := Request{
		Method:      method,
		URI:         path,
		Headers:     header,
		QueryParams: params,
	}

	return &request
}

func (r *Request) GetPathValue(key string) string {

	value, ok := r.PathParams[key]

	if ok {
		return value
	}

	return ""
}
