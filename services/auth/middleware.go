package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nirav114/url-shortner-backend.git/config"
	"github.com/nirav114/url-shortner-backend.git/types"
)

type contextKey string

const userContextKey contextKey = "user"

type UserClaims struct {
	UserID    int64 `json:"userID"`
	ExpiredAt int64 `json:"exp"`
}

func JWTMiddleware(store types.UserStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var secret = []byte(config.EnvConfig.JWTSecret)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userClaims := &UserClaims{}

			if userID, ok := claims["userid"].(string); ok {
				userClaims.UserID, err = strconv.ParseInt(userID, 10, 64)
				if err != nil {
					http.Error(w, "Invalid userID in token", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Invalid userID in token", http.StatusUnauthorized)
				return
			}

			user, err := store.GetUserByID(userClaims.UserID)
			if err != nil || user == nil {
				http.Error(w, "User does not exist", http.StatusUnauthorized)
				return
			}

			if exp, ok := claims["expiredAt"].(float64); ok {
				userClaims.ExpiredAt = int64(exp)
			} else {
				http.Error(w, "Invalid expiration time in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, userClaims)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}

func GetUserClaimsFromContext(ctx context.Context) (*UserClaims, bool) {
	userClaims, ok := ctx.Value(userContextKey).(*UserClaims)
	return userClaims, ok
}
