package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/debug/info", handleServiceInfo)

	fmt.Println("Server started: 8090")
	http.ListenAndServe(":8090", nil)
}

func handleServiceInfo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
