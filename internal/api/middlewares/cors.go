package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrgins = []string{
	"https://localhost:3000",
	"https://www.frontend.com",
}

// Here Cors stand for corss Origin resource sharing which helps to understand what origins are allowed to access the api here we will use an exmaple that our Api can be accesed by two domains
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println(origin)

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Not allowed by the CORS", http.StatusForbidden)
			return
		}

		w.Header().Set("Acess-Control-Allow-Headers", "Content-Type, Authorization")   //What headers can be used in the request
		w.Header().Set("Acess-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE") //What methods can be used in the request
		w.Header().Set("Access-Control-Allow-Credentials", "true")                     //tells browsers whether the server allows credentials to be included in cross-origin HTTP requests.
		w.Header().Set("Acess-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrgins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
