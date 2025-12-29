package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello From The Root Route"))
	fmt.Println("Hello From The Root Route")
}
