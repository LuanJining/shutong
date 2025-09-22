package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Claims represent JWT claims used by the platform.
type Claims struct {
    UserID string `json:"user_id"`
    Phone  string `json:"phone"`
    jwt.RegisteredClaims
}

// Manager handles JWT generation and validation.
type Manager struct {
    signingKey []byte
    ttl        time.Duration
}

// NewManager creates a JWT manager with the given secret and TTL in seconds.
func NewManager(signingKey string, ttlSeconds int) *Manager {
    if ttlSeconds <= 0 {
        ttlSeconds = 3600
    }
    return &Manager{
        signingKey: []byte(signingKey),
        ttl:        time.Duration(ttlSeconds) * time.Second,
    }
}

// Generate creates a signed JWT for the provided identity.
func (m *Manager) Generate(userID, phone string) (string, error) {
    expiresAt := time.Now().Add(m.ttl)
    claims := Claims{
        UserID: userID,
        Phone:  phone,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   userID,
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(m.signingKey)
}

// Validate parses and validates a JWT returning claims on success.
func (m *Manager) Validate(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return m.signingKey, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token claims")
}
