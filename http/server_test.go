package http

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestRegisterValidURI(t *testing.T) {

	s := New()

	if !testRegisteringRoute(t, s, "/test", "/test") {
		return
	}

	if !testRegisteringRoute(t, s, "/user/{id}/email/{gmail}", "/user/*/email/*") {
		return
	}

}

func TestRegisterInvalidURI(t *testing.T) {

	s := New()

	err := s.Reg("test", func(r Request) {
		fmt.Println("Should not run")
	})

	if err == nil {
		t.Fatalf("err should NOT be nil. Expected=Error")
	}

	if s.routes["test"].handler != nil {
		t.Fatalf("s.routes['test'] should not be set. Expected nil value")
	}
}

func testRegisteringRoute(t *testing.T, s *HTTPServer, routeUri, expectedRoute string) bool {

	err := s.Reg(routeUri, func(r Request) {
	})

	if err != nil {
		t.Errorf("err is not nil. Got=%s", err.Error())
		return false
	}

	if s.routes[expectedRoute].handler == nil {
		t.Errorf("s.routes['%s'].handler is  nil. tried to register=%s, expected=%s", expectedRoute, routeUri, expectedRoute)
		return false
	}

	uriSl := strings.Split(routeUri, "/")
	if !slices.Equal(uriSl, s.routes[expectedRoute].uri) {
		t.Errorf("s.routes['%s'].uri(%v) not equal to routeUri(%v)", expectedRoute, uriSl, s.routes[expectedRoute].uri)
		return false
	}

	return true

}
