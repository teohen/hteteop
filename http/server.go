package http

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type RouteHandler func(Request)

type RouteRegister struct {
	handler RouteHandler
	uri     []string
}

type HTTPServer struct {
	port   int
	routes map[string]RouteRegister
}

func New() *HTTPServer {
	s := HTTPServer{routes: make(map[string]RouteRegister)}

	return &s
}

func (s *HTTPServer) Listen(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Error:", err)
		return errors.New("Unable to listen to port: " + err.Error())
	}
	defer listener.Close()

	fmt.Println("Server listening to PORT: ", port)
	s.port = port

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *HTTPServer) handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	reqstr := ""

	n := 1024
	var err error

	for {
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		reqstr += fmt.Sprintf("%s", buffer[:n])
		if n < 1024 {
			break
		}
	}

	parsedReq := *ParseRequest(reqstr)

	maskedRoute := getRegisteredRoute(parsedReq.URI, s.routes)
	if maskedRoute != "" {
		routeRegister := s.routes[maskedRoute]
		parsedReq.ParsePathParams(routeRegister.uri)
		go routeRegister.handler(parsedReq)
	} else {
		// TODO: implement the 404 response later
		fmt.Println("404 Not Found")
		// return 404
	}
}

func getRegisteredRoute(uri string, routes map[string]RouteRegister) string {
	uriSl := strings.Split(uri, "/")

	for _, route := range routes {
		if len(uriSl) == len(route.uri) {
			match := true
			for i, part := range route.uri {
				if part != "" && part[0] != '{' && part != uriSl[i] {
					match = false
				}
			}
			if match {
				return getMaskedURI(route.uri)
			}
		}
	}

	return ""
}

func getMaskedURI(uriSlice []string) string {
	res := ""
	for _, part := range uriSlice {
		if part != "" {
			if part[0] != '{' {
				res += part
			} else {
				res += "*"
			}
		}
		res += "/"
	}
	return res[:len(res)-1]
}

func (s *HTTPServer) Reg(uri string, handler RouteHandler) error {

	if uri[0] != '/' {
		return errors.New("uri should start with / character")
	}

	rr := RouteRegister{
		handler: handler,
		uri:     strings.Split(uri, "/"),
	}

	maskedUri := getMaskedURI(strings.Split(uri, "/"))
	s.routes[maskedUri] = rr

	return nil
}
