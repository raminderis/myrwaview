package main

import (
	"fmt"
	"net/http"
)

func thisIsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome website!</h1>")
}

func main() {
	http.HandleFunc("/", thisIsHandler)
	http.ListenAndServe(":3000", nil)
}
