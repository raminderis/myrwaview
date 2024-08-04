package main

import (
	"fmt"
	"net/http"
)

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p> To get in touch with me, email me at: <a href='mailto:rami@live.com'>rami@live.com</a></p>")

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Home Page</h1><p> Here we will have some fun with Real World Assets from On Chain</p>")

}

type Router struct {
}

func (x Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, "Page Not Found", http.StatusNotFound)
	}
}

// func pathHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		homeHandler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)
// 	default:
// 		http.Error(w, "Page Not Found", http.StatusNotFound)
// 	}
// }

func main() {
	var router Router
	http.ListenAndServe(":3000", router)
}
