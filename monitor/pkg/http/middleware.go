package http

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type Role string

const (
	RolePassenger Role = "passenger"
	RoleDriver    Role = "driver"
)

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role Role   `json:"role"`
}

func VerifyFakeToken(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			_ = ErrorJSON(w, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 3)
		if len(parts) != 2 || parts[0] != "Bearer" {
			_ = ErrorJSON(w, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid authorization header",
			})
			return
		}

		token := parts[1]
		profile, err := extractFakeProfile(token)
		if err != nil {
			_ = ErrorJSON(w, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "failed to extract profile",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "profile", profile)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractFakeProfile(token string) (*Profile, error) {
	if strings.HasPrefix(token, string(RoleDriver)) {
		return &Profile{
			ID:   token,
			Name: "Driver " + token,
			Role: RoleDriver,
		}, nil
	}

	if strings.HasPrefix(token, string(RolePassenger)) {
		return &Profile{
			ID:   token,
			Name: "Passenger " + token,
			Role: RolePassenger,
		}, nil
	}

	return nil, errors.New("invalid token")
}
