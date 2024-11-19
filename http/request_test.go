package http

import (
	"testing"
)

func TestParseRequest(t *testing.T) {

	r := "GET /users HTTP/1.1\r\n"
	r += "User-Agent: curl/7.54.1\r\n"
	r += "Accept: */*\r\n"
	r += "Host: localhost:8080\r\n"
	r += "Accept-Encoding: gzip, deflate, br\r\n"
	r += "Connection: keep-alive\r\n"

	req := ParseRequest(r)

	if req.URI != "/users" {
		t.Fatalf("req.URI has wrong value. want=/users, got=%s", req.URI)
	}

	if req.Method != "GET" {
		t.Fatalf("req.Method has wrong value. want=GET, got=%s", req.Method)
	}

	if len(req.Headers) != 5 {
		t.Fatalf("req.Header has wrong length. want=5, got=%d", len(req.Headers))
	}

	if !testHeader(t, req, "User-Agent", "curl/7.54.1") {
		return
	}

	if !testHeader(t, req, "Accept", "*/*") {
		return
	}

	if !testHeader(t, req, "Host", "localhost:8080") {
		return
	}

	if !testHeader(t, req, "Accept-Encoding", "gzip, deflate, br") {
		return
	}

	if !testHeader(t, req, "Connection", "keep-alive") {
		return
	}
}
func testHeader(t *testing.T, r *Request, expectedKey, expectedValue string) bool {
	value, ok := r.Headers[expectedKey]

	if !ok {
		t.Errorf("r.Headers[expectedKey] does not exist.")
		return false
	}

	if value != expectedValue {
		t.Errorf("r.Header[%s] has the wrong value. want='%s' | got='%s'", expectedKey, expectedValue, value)
		return false
	}

	return true
}
