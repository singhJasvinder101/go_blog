// this package allows you to initialize the jwt utility
// once in main function with secret field and no need
// to pass this modules netsted way (like prop drilling)
// internal pointer singleton instance(package level) and
// methods supporting exported methods

package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTUtils struct {
	Secret string
}

// package-level singleton
var jwtInstance *JWTUtils

func Init(secret string) *JWTUtils {
	if jwtInstance == nil {
		jwtInstance = &JWTUtils{Secret: secret}
	}
	return jwtInstance
}

// package-level functions (for using anywhere)
func GenerateToken(userId int, name, email string) (string, error) {
	if jwtInstance == nil {
		return "", fmt.Errorf("JWT not initialized")
	}

	return jwtInstance.generateToken(userId, name, email)
}

func VerifyToken(token string) (*Claims, error) {
	if jwtInstance == nil {
		return nil, fmt.Errorf("JWT not initialized")
	}

	return jwtInstance.verifyToken(token)
}

// internal pointer receiver methods
func (j *JWTUtils) generateToken(userId int, name, email string) (string, error) {
	claims := Claims{
		UserID: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go_blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Println(token)

	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", fmt.Errorf("error when signing token %w", err)
	}
	return signedToken, nil
}

func (j *JWTUtils) verifyToken(tokenStr string) (*Claims, error) {
	if jwtInstance == nil {
		return nil, fmt.Errorf("JWT not initialized")
	}

	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// allow only signingHS256
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(j.Secret), nil
	})
	if err != nil || !parsedToken.Valid {
		return &Claims{}, fmt.Errorf("wrong credentials")
	}

	claims := parsedToken.Claims.(Claims)
	return &claims, nil
}
