package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	mw "schoolREST/internal/api/middlewares"
	"schoolREST/internal/api/router"
)

func main() {

	port := ":3000"

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	//These will  be used again but for now they are commented out!!!

	// rl := mw.NewRateLimiter(8, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	/* 	secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.Security_headers(mw.Compression(mw.Hpp(hppOptions)(mux)))))) //This specific middleware chaining reflects that the request travels from Cors to Hpp and then request travels from Hpp to cors again!!
	 */

	// secureMux := utils.ApplyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.Security_headers, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)
	//For designing the end points we will be using only security headers since we have tested the middlewares we will test them again after endpoints have been made

	router := router.Router()
	secureMux := mw.Security_headers(router)

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server starting on port ", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln(err.Error())
	}

}
