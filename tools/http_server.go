package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func payload(w http.ResponseWriter, req *http.Request) {
	fmt.Println("---------------------- request ----------------------")
	responseData, _ := ioutil.ReadAll(req.Body)
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Printf("%v: %v\n", name, h)
		}
	}
	fmt.Printf("\n%s\n\n", responseData)
	fmt.Fprintf(w, "%s", responseData)
}

func main() {
	http.HandleFunc("/payload", payload)
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
