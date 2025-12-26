package middlewares

import "net/http"

func Security_headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off") //This doesn't allow the browser to prefetch DNS for domains while users are browsing Helps reduce DNS related attacks

		w.Header().Set("X-Frame-Options", "DENY")                                                   //This helps so that the webpage cannot be displayed inside Iframe to prevent against clickjacking attacks
		w.Header().Set("X-XSS-Protection", "1; mode=block")                                         //Modern web browsers contain protection against XSS aka cross site scripting against but is not sometines enabled by default we have enabled it here
		w.Header().Set("X-Content-Type-Options", "nosniff")                                         //Server the actual content type Trust the content type provided by the API
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload") //Only connect via HTTPS in the provided max age in seconds
		w.Header().Set("Content-Security-Policy", "default-src 'self'")                             //Only load resources from the safe resources or same resources
		w.Header().Set("Referred-Policy", "no-referrer")                                            //No referer information is sent
		w.Header().Set("X-Powered-By", "Django")                                                    //You can throw the attacker by making them think you are using the Django framework for backend
		w.Header().Set("Server", "")                                                                //Tells the browser what Software is being used in the server currently it is empty due to us using none of those
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")                                 //defines a meta-policy that controls whether site resources can be accessed cross-origin by a document running in a web client
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")           // Contains instruction on how to control caching in browsers
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")                               //Can a browser bloack resources from different origins
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")                                 // allows a website to control whether a new top-level document, opened using Window.open() or by navigating to a new page
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")                              //response header configures the current document's policy for loading and embedding cross-origin resources.
		w.Header().Set("Acess-Control-Allow-Headers", "Content-Type, Authorization")                //What headers can be used in the request
		w.Header().Set("Acess-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")              //What methods can be used in the request
		w.Header().Set("Access-Control-Allow-Credentials", "true")                                  //tells browsers whether the server allows credentials to be included in cross-origin HTTP requests.
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")                   //provides a mechanism to allow and deny the use of browser features in a document or within any <iframe> elements in the document
		next.ServeHTTP(w, r)
	})
}

//Basic Middleware skeleton
// func security_headers(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//next.ServeHTTP(w,r)
// 	})
// }
