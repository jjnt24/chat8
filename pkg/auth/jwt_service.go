package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	CityID   int64  `json:"city_id"`
	jwt.RegisteredClaims
}

var ErrTokenExpired = errors.New("token expired")

func (s JWTService) NewAccess(userID int64, username string, cityID int64) (string, error) {
	claims := Claims{UserID: userID, Username: username, CityID: cityID, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.AccessTTL))}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.AccessSecret)
}

func (s JWTService) NewRefresh(userID int64, username string, cityID int64) (string, error) {
	claims := Claims{UserID: userID, Username: username, CityID: cityID, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.RefreshTTL))}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.RefreshSecret)
}

func (s JWTService) ParseAccess(tok string) (*Claims, error) {
	var c Claims

	t, err := jwt.ParseWithClaims(tok, &c, func(t *jwt.Token) (any, error) {
		return s.AccessSecret, nil
	})
	if err != nil {
		// === JWT v5: cek apakah error karena token expired ===
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, err
	}

	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	return &c, nil
}

func (s JWTService) ParseRefresh(tok string) (*Claims, error) {
	var c Claims
	t, err := jwt.ParseWithClaims(tok, &c, func(t *jwt.Token) (any, error) { return s.RefreshSecret, nil })
	if err != nil || t == nil || !t.Valid {
		return nil, err
	}
	return &c, nil
}
