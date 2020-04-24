package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tatrasoft/fyp-rest-backend-service/utils/auth"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/response_models"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// getting the token from the header
		tokenFromHeader := request.Header.Get("Authorization")

		// checks if there is auth token present in the request header
		if tokenFromHeader == "" {
			response.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(response).Encode(response_models.ErrorResponse{
				ErrorCode: http.StatusUnauthorized,
				ErrorMessage:  "Missing auth token",
			})

			return
		}

		// retrieves the jwt token from the header
		splitToken := strings.Split(tokenFromHeader, "Bearer")
		tokenFromHeader = strings.TrimSpace(splitToken[1])

		// check if the token is valid
		isValid, _, err := auth.IsValidJWTToken(tokenFromHeader)
		if err != nil && !isValid {
			response.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(response).Encode(response_models.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: fmt.Sprintf("invalid jwt token: %v", err),
			})

			return
		}

		next.ServeHTTP(response, request)
	})
}
