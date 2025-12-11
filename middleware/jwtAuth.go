package middleware

import (
	"net/http"
	"strings"
	"vetclinic-rest-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = utils.JwtSecret

// JWTAuth middleware for validation token + role
func JWTAuth(requiredRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // fetch header Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // parse token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        role, ok := claims["role"].(string)
        if !ok {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Check if the role exactly the one that needed
        allowed := false
        for _, r := range requiredRoles {
            if role == r {
                allowed = true
                break
            }
        }

        if !allowed {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }

        // Extract user_id from claims
        userId, ok := claims["user_id"].(string)
        if !ok {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Save role and user_id to context
        c.Set("role", role)
        c.Set("user_id", userId)

        c.Next()
    }
}