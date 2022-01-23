package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/jiny0x01/squareheart_backend/api/auth"
)

func main() {
	fmt.Println("Server start running at 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/signup", auth.SignUp)

	http.ListenAndServe(":8080", nil)
}
