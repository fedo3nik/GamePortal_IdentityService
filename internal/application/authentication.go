package application

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	c "github.com/fedo3nik/GamePortal_IdentityService/config"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/domain/entities"
	e "github.com/fedo3nik/GamePortal_IdentityService/internal/error"
)

const (
	aKey            = "access"
	rKey            = "refresh"
	timeDurationSec = 300
)

type AuthService struct {
	config *c.Config
}

type RefreshTokenCustomClaims struct {
	UserID    string
	CustomKey string
	KeyType   string
	Claims    jwt.StandardClaims
}

type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	Claims  jwt.StandardClaims
}

func (a *AccessTokenCustomClaims) Valid() error {
	if a.KeyType == aKey {
		return nil
	}

	return e.ErrJwtClaims
}

func (r *RefreshTokenCustomClaims) Valid() error {
	if r.KeyType == rKey {
		return nil
	}

	return e.ErrJwtClaims
}

func GenerateCustomKey(userID, tokenHash string) (string, error) {
	h := hmac.New(sha256.New, []byte(tokenHash))

	_, err := h.Write([]byte(userID))
	if err != nil {
		return "", err
	}

	sha := hex.EncodeToString(h.Sum(nil))

	return sha, nil
}

func (a *AuthService) GenerateRefreshToken(user *entities.User) (string, error) {
	customKey, err := GenerateCustomKey(user.ID, user.TokenHash)
	if err != nil {
		return "", err
	}

	tokenType := rKey

	claims := RefreshTokenCustomClaims{
		UserID:    user.ID,
		CustomKey: customKey,
		KeyType:   tokenType,
		Claims: jwt.StandardClaims{
			Issuer: "identityservice.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(a.config.RefreshPrivateKey)
	if err != nil {
		return "", err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims)

	return token.SignedString(signKey)
}

func (a *AuthService) GenerateAccessToken(user *entities.User) (string, error) {
	tokenType := aKey

	claims := AccessTokenCustomClaims{
		UserID:  user.ID,
		KeyType: tokenType,
		Claims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(timeDurationSec)).Unix(),
			Issuer:    "identityservice.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(a.config.AccessPrivateKey)
	if err != nil {
		return "", err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims)

	return token.SignedString(signKey)
}

func NewAuthService(config *c.Config) *AuthService {
	return &AuthService{config: config}
}
