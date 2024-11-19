package main

import (
	"fmt"
	"log"

	"github.com/teohen/hteteop/http"
)

func handler(r http.Request) {
	fmt.Println("Request to: ", r.URI)
	fmt.Println("Request method: ", r.Method)
	fmt.Println("Request Headers: ", r.Headers)
}

func main() {
	s := http.New()

	s.Reg("/users", handler)

	err := s.Listen(8080)

	if err != nil {
		log.Fatal("error listening. ", err.Error())
	}
}
