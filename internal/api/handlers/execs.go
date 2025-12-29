package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("This is the GET Method for Executives routes"))
	case http.MethodPost:
		w.Write([]byte("This is the PUT Method for Executives routes"))
	case http.MethodDelete:
		w.Write([]byte("This is the DELETE Method for Executives routes"))
	case http.MethodPatch:
		w.Write([]byte("This is the PATCH Method for Executives routes"))
	}
}
