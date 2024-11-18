package http

import (
	"testing"
)

func TestParseRequest(t *testing.T) {

	r := "GET /users HTTP/1.1\r\nUser-Agent: curl/7.54.1\r\nAccept: */*\r\nHost: localhost:8080\r\nAccept-Encoding: gzip, deflate, br\r\nConnection: keep-alive\r\n"

	req := ParseRequest(r)

	if req.URI != "/users" {
		t.Fatalf("req.URI has wrong value. want=/users, got=%s", req.URI)
	}

	if req.Method != "GET" {
		t.Fatalf("req.Method has wrong value. want=GET, got=%s", req.Method)
	}

}
