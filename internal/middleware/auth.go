package middleware

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if len(header) < 8 || header[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется токен"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(header[7:], func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "некорректный токен"})
			c.Abort()
			return
		}

		rawUserID, ok := claims["user_id"].(float64)
		if !ok || math.IsNaN(rawUserID) || math.IsInf(rawUserID, 0) ||
			rawUserID <= 0 || math.Trunc(rawUserID) != rawUserID ||
			rawUserID >= math.Exp2(64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "некорректный пользователь в токене"})
			c.Abort()
			return
		}

		c.Set("user_id", uint(rawUserID))
		c.Next()
	}
}
