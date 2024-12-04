package http

import (
	"testing"
)

func TestParseGETRequestWithQueryParams(t *testing.T) {

	r := "GET /users?name=te&email=gmail HTTP/1.1\r\n"
	r += "User-Agent: curl/7.54.1\r\n"
	r += "Accept: */*\r\n"
	r += "Host: localhost:8080\r\n"
	r += "Accept-Encoding: gzip, deflate, br\r\n"
	r += "Connection: keep-alive\r\n"

	req := ParseRequest(r)

	if !testRequestLine(t, req, "/users", "GET", 5) {
		return
	}

	if !testRequestMapProperty(t, req.Headers, "Headers", "User-Agent", "curl/7.54.1") {
		return
	}

	if !testRequestMapProperty(t, req.Headers, "Headers", "Accept", "*/*") {
		return
	}

	if !testRequestMapProperty(t, req.Headers, "Headers", "Host", "localhost:8080") {
		return
	}

	if !testRequestMapProperty(t, req.Headers, "Headers", "Accept-Encoding", "gzip, deflate, br") {
		return
	}

	if !testRequestMapProperty(t, req.Headers, "Headers", "Connection", "keep-alive") {
		return
	}

	if !testRequestMapProperty(t, req.QueryParams, "Params", "name", "te") {
		return
	}
	if !testRequestMapProperty(t, req.QueryParams, "Params", "email", "gmail") {
		return
	}

}

func testRequestLine(t *testing.T, r *Request, expUri, expMethod string, expheaderLength int) bool {

	if r.URI != expUri {
		t.Fatalf("Request.URI has wrong value. want=%s, got=%s", expUri, r.URI)
	}

	if r.Method != expMethod {
		t.Fatalf("Request.Method has wrong value. want=%s, got=%s", expMethod, r.Method)
	}

	if len(r.Headers) != expheaderLength {
		t.Fatalf("req.Header has wrong length. want=%d, got=%d", expheaderLength, len(r.Headers))
	}

	return true

}

func testRequestMapProperty(t *testing.T, mapped map[string]string, property, expectedKey, expectedValue string) bool {

	if mapped == nil {
		t.Errorf("Request.%s does not exist.", property)
		return false
	}

	value, ok := mapped[expectedKey]

	if !ok {
		t.Errorf("Request.%s['%s'] does not exist. Expected=%s, got=nil", property, expectedKey, expectedKey)
	}

	if value != expectedValue {
		t.Errorf("Request.%s['%s'] has the wrong value. want='%s' | got='%s'", property, expectedKey, expectedValue, value)
		return false
	}

	return true
}
