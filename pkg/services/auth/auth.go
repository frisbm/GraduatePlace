package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/frisbm/graduateplace/pkg/models/user"
)

type AuthService struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewAuthService(secretKey string) *AuthService {
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	return &AuthService{
		tokenAuth: tokenAuth,
	}
}

func (a *AuthService) GetTokenAuth() *jwtauth.JWTAuth {
	return a.tokenAuth
}

func (a *AuthService) generateAccessToken(uuid string) (string, error) {
	tokenMap := make(map[string]interface{})
	tokenMap["uuid"] = uuid
	tokenMap["kind"] = "access"
	jwtauth.SetIssuedNow(tokenMap)
	jwtauth.SetExpiryIn(tokenMap, time.Hour)
	_, tokenString, err := a.tokenAuth.Encode(tokenMap)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthService) generateRefreshToken(uuid string) (string, error) {
	tokenMap := make(map[string]interface{})
	tokenMap["uuid"] = uuid
	tokenMap["kind"] = "refresh"

	jwtauth.SetIssuedNow(tokenMap)
	jwtauth.SetExpiryIn(tokenMap, time.Hour*24*30)
	_, tokenString, err := a.tokenAuth.Encode(tokenMap)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthService) GenerateTokens(uuid string) (*user.AuthTokens, error) {
	accessToken, err := a.generateAccessToken(uuid)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateRefreshToken(uuid)
	if err != nil {
		return nil, err
	}
	return &user.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

var invalidRefreshTokenErr = errors.New("invalid refresh token")
var invalidAccessTokenErr = errors.New("invalid access token")
var accessTokenExpired = errors.New("access token expired")

func (a *AuthService) ParseRefreshToken(refreshToken string) (*uuid.UUID, error) {
	token, err := a.tokenAuth.Decode(refreshToken)
	if err != nil {
		return nil, err
	}

	kind, ok := token.Get("kind")
	if !ok {
		return nil, invalidRefreshTokenErr
	}

	if fmt.Sprintf("%v", kind) != "refresh" {
		return nil, invalidRefreshTokenErr
	}

	if token.Expiration().UTC().Unix() < time.Now().UTC().Unix() {
		return nil, invalidRefreshTokenErr
	}

	uuidRes, ok := token.Get("uuid")
	if !ok {
		return nil, invalidRefreshTokenErr
	}

	uuidStr, ok := uuidRes.(string)
	if !ok {
		return nil, invalidRefreshTokenErr
	}

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, invalidRefreshTokenErr
	}
	return &uuid, nil
}

func (a *AuthService) ParseAccessToken(accessToken string) (*uuid.UUID, error) {
	token, err := a.tokenAuth.Decode(accessToken)
	if err != nil {
		return nil, err
	}

	kind, ok := token.Get("kind")
	if !ok {
		return nil, invalidAccessTokenErr
	}

	if fmt.Sprintf("%v", kind) != "access" {
		return nil, invalidAccessTokenErr
	}

	if token.Expiration().UTC().Unix() < time.Now().UTC().Unix() {
		return nil, accessTokenExpired
	}

	uuidRes, ok := token.Get("uuid")
	if !ok {
		return nil, invalidAccessTokenErr
	}

	uuidStr, ok := uuidRes.(string)
	if !ok {
		return nil, invalidAccessTokenErr
	}

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, invalidAccessTokenErr
	}
	return &uuid, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePasswordCorrect(hashedPassword, submittedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(submittedPassword))
}
