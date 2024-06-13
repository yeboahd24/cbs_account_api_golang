package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"strings"
	// "io"
)

const (
	authorizationHeaderKey = "Authorization"
)

// content-type middleware
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(authorizationHeaderKey)
		tokenRequest := c.GetHeader(authorizationHeaderKey)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "No token provided", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Token format is invalid", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]
		// fmt.Println("Token string:", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secret := []byte(os.Getenv("SECRET"))
			// fmt.Println("Using secret:", string(secret))
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// fmt.Printf("Token claims: %+v\n", token.Claims)
			return secret, nil
		})

		if err != nil {
			fmt.Println("Token parsing error:", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error(), StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// fmt.Printf("Token is valid. Claims: %+v\n", claims)
			payload := map[string]string{"token": tokenString}
			// fmt.Println("payload", payload)
			jsonPayload, err := json.Marshal(payload)
			// fmt.Println("jsonPayload", jsonPayload)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error encoding payload", StatusCode: http.StatusInternalServerError})
				c.Abort()
				return
			}

			authURL := os.Getenv("AUTH_URL")
			request, err := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(jsonPayload))
			request.Header.Set("Authorization", tokenRequest)

			client := &http.Client{}
			response, err := client.Do(request)
			// fmt.Println("response", response)

			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error contacting auth service", StatusCode: http.StatusInternalServerError})
				c.Abort()
				return
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token", StatusCode: http.StatusUnauthorized})
				c.Abort()
				return
			}

			var result map[string]interface{}
			err = json.NewDecoder(response.Body).Decode(&result)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error decoding response", StatusCode: http.StatusInternalServerError})
				c.Abort()
				return
			}

			c.Set("userID", result["user_id"])
			c.Next()
		} else {
			fmt.Println("Token claims invalid or token is not valid")
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}
	}
}
