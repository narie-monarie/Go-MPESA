package main

import (
	stkpush "narie-monarie/stk-push"
	"net/http"
)

func main() {
	http.HandleFunc("POST /", stkpush.GetSTKPush)
	http.ListenAndServe(":8080", nil)
}
