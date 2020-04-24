package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/auth"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/response_models"
	"net/http"
	"strings"
)

// TODO implement admin middleware
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// getting the token from the header
		tokenFromHeader := request.Header.Get("Authorization")

		// checks if there is auth token present in the request header
		if tokenFromHeader == "" {
			response.WriteHeader(http.StatusForbidden)
			json.NewEncoder(response).Encode(response_models.ErrorResponse{
				ErrorCode: http.StatusForbidden,
				ErrorMessage:  "Missing auth token",
			})

			return
		}

		// retrieves the jwt token from the header
		splitToken := strings.Split(tokenFromHeader, "Bearer")
		tokenFromHeader = strings.TrimSpace(splitToken[1])

		// check if the token is valid
		_, isAdmin, err := auth.IsValidJWTToken(tokenFromHeader)
		if err != nil && !isAdmin {
			response.WriteHeader(http.StatusForbidden)
			json.NewEncoder(response).Encode(response_models.ErrorResponse{
				ErrorCode:    http.StatusForbidden,
				ErrorMessage: fmt.Sprintf("missing admin privilages: %v", err),
			})

			return
		}

		next.ServeHTTP(response, request)
	})
}
