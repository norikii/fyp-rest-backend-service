package middleware

import "net/http"


//w.Header().Add("content-type", "application/json")
//w.Header().Add("Access-Control-Allow-Origin", "*")
//w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
//w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
//w.Header().Add("Access-Control-Max-Age", "86400")

// CommonMiddleware sets the content type and access control headers
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		//if request.Method == "OPTIONS" {
		//	response.Header().Set("Access-Control-Allow-Origin", "*")
		//	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		//	response.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
		//	response.Header().Set("Content-Type", "application/json")
		//	return
		//}
		response.Header().Add("Content-Type", "application/json")
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(response, request)
	})
}
