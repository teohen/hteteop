package main

import (
	"fmt"
	"log"

	"github.com/teohen/hteteop/http"
)

func handler(r http.Request) {
	fmt.Println("path value sent to route: ",
		r.GetPathValue("id"))
}

func main() {
	s := http.New()

	s.Reg("/users/{id}", handler)

	err := s.Listen(8080)

	if err != nil {
		log.Fatal("error listening. ", err.Error())
	}
}
