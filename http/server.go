package http

import (
	"errors"
	"fmt"
	"net"
)

type RouteHandler func(Request)

type HTTPServer struct {
	port   int
	routes map[string]RouteHandler
}

func New() *HTTPServer {
	s := HTTPServer{routes: make(map[string]RouteHandler)}

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
	go s.routes[parsedReq.URI](parsedReq)
}

func (s *HTTPServer) Reg(uri string, hand RouteHandler) {
	s.routes[uri] = hand
}
