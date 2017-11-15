package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//fmt.Fprint(w, r.URL.Path)
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>GO, GO, GO, GOLANG!</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "Please contact us at <a href=\"mailto:hihelloreid@gmail.com\">hihelloreid</a>.")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
