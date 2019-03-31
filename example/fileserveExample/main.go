package main

import (
	"net/http"
	"github.com/rbxb/fileserve"
)

func main() {
	server := fileserve.NewServer("./root", nil)
	http.ListenAndServe(":8080", server)
}