package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
		for name, value := range r.Header {
			fmt.Printf("%s %s\n", name, value)
		}
		fmt.Println()
		if buf, err := io.ReadAll(r.Body); err != nil {
			fmt.Printf("Failed to read request body: %v", err)
		} else {
			fmt.Printf("%s\n", string(buf))
		}
		fmt.Println()
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
