package main

import (
	"fmt"
	"net/http"
)

var (
	Commit    string
	BuildTime string
)

func main() {
	http.HandleFunc("/debug/info", handleServiceInfo)

	fmt.Println("Server started: 8090")
	http.ListenAndServe(":8090", nil)
}

func handleServiceInfo(w http.ResponseWriter, req *http.Request) {
	format := req.Header.Get("Accept")
	var resp string

	switch format {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		resp = fmt.Sprintf(`{"commit":"%s","date":"%s"}`, Commit, BuildTime)
	default:
		resp = "Commit SHA1: " + Commit + " " + BuildTime
	}
	fmt.Fprintln(w, resp)
}
