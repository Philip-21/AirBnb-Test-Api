package middleware

import (
	"airbnb/models"
	"airbnb/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtUserClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

type JwtPropertyClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateUserToken(ID uuid.UUID) (string, error) {
	claims := &JwtUserClaims{
		ID: ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", fmt.Errorf("unable to generate token %s", err)
	}
	return t, nil
}

func AuthUser(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		} else {
			tokenString = strings.TrimSpace(authHeader)
		}

		if tokenString == "" {
			log.Println("Authorization header invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JwtUserClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		user, err := userRepo.GetUserByID(c.Request.Context(), claims.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin not found"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func GetUser(ctx *gin.Context) (*models.User, error) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}
	User, ok := user.(*models.User)
	if !ok {
		return nil, fmt.Errorf("invalid user type in context")
	}
	return User, nil
}

func GeneratePropertyOwnerToken(ID uuid.UUID) (string, error) {
	claims := &JwtPropertyClaims{
		ID: ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", fmt.Errorf("unable to generate token %s", err)
	}
	return t, nil
}

func AuthPropertyOwner(propertyRepo *repository.PropertyRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		} else {
			tokenString = strings.TrimSpace(authHeader)
		}

		if tokenString == "" {
			log.Println("Authorization header invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &JwtPropertyClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JwtPropertyClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		property, err := propertyRepo.GetPropertyOwnerByID(c.Request.Context(), claims.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "property not found"})
			c.Abort()
			return
		}
		c.Set("owner", property)
		c.Next()
	}
}

func GetPropertyOwner(ctx *gin.Context) (*models.PropertyOwner, error) {
	propertyOwner, exists := ctx.Get("owner")
	if !exists {
		return nil, fmt.Errorf("owner not found in context")
	}
	owner, ok := propertyOwner.(*models.PropertyOwner)
	if !ok {
		return nil, fmt.Errorf("invalid owner type in context")
	}
	return owner, nil
}
