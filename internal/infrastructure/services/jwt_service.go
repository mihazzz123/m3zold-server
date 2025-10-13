package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mihazzz123/m3zold-server/internal/domain/auth"
	"github.com/mihazzz123/m3zold-server/internal/domain/services"
)

type JWTService struct {
	secretKey []byte
	issuer    string
	audience  string
}

func NewJWTService(secretKey, issuer, audience string) services.JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		audience:  audience,
	}
}

func (s *JWTService) GenerateToken(userID, email, userName string) (string, error) {
	claims := auth.Claims{
		UserID:   userID,
		Email:    email,
		UserName: userName,
		Exp:      time.Now().Add(15 * time.Minute),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   claims.UserID,
		"email":     claims.Email,
		"user_name": claims.UserName,
		"exp":       claims.Exp.Unix(),
		"iss":       s.issuer,
		"aud":       s.audience,
		"iat":       time.Now().Unix(),
	})

	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем issuer и audience
		if iss, ok := claims["iss"].(string); !ok || iss != s.issuer {
			return nil, fmt.Errorf("invalid issuer")
		}

		if aud, ok := claims["aud"].(string); !ok || aud != s.audience {
			return nil, fmt.Errorf("invalid audience")
		}

		userID, _ := claims["user_id"].(string)
		email, _ := claims["email"].(string)
		userName, _ := claims["user_name"].(string)
		exp, _ := claims["exp"].(float64)

		return &auth.Claims{
			UserID:   userID,
			Email:    email,
			UserName: userName,
			Exp:      time.Unix(int64(exp), 0),
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
