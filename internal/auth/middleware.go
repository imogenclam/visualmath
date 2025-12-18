package auth

import (
    "context"
    "net/http"
   //"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
    UserID   int    `json:"user_id"`
    Login    string `json:"login"`
    UserType string `json:"user_type"`
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userClaims := UserClaims{
            UserID:   1,
            Login:    "test_user",
            UserType: "teacher", // Или "student", "admin"
        }
        
        ctx := context.WithValue(r.Context(), UserContextKey, userClaims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            next.ServeHTTP(w, r)
        })
    }
}

func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
    user, ok := ctx.Value(UserContextKey).(UserClaims)
    return &user, ok
}