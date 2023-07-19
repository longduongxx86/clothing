package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		tokenString := strings.Split(authHeader, " ")[1]
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Token")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
			}
			return []byte("clothing_token"), nil // Thay "secret" bằng một khóa ngẫu nhiên an toàn
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("email", claims["email"])
			c.Set("fullname", claims["fullname"])
			c.Set("role_id", claims["role_id"])
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
	}
}
