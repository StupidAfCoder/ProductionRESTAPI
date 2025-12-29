package handlers

import (
	"fmt"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("This is the GET Method for students routes"))
	case http.MethodPost:
		w.Write([]byte("This is the PUT Method for students routes"))
	case http.MethodDelete:
		w.Write([]byte("This is the DELETE Method for students routes"))
	case http.MethodPatch:
		w.Write([]byte("This is the PATCH Method for students routes"))
	}

}
