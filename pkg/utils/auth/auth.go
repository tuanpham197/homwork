package auth_util

import (
	commonapp "Ronin/common"
	"Ronin/internal/auth/services/request"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

var (
	SECRET_KEY = os.Getenv("SECRET_KEY")
)

type Keys struct {
	PublicKey  string
	PrivateKey string
}

type Permission struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`
	GuardName string     `json:"guard_name"`
	Role      []Role     `json:"roles" gorm:"many2many:role_has_permissions;"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Role struct {
	ID          uint                `json:"id" gorm:"primaryKey"`
	Name        string              `json:"name"`
	GuardName   string              `json:"guard_name" mapstructure:"guard_name"`
	Permissions []entity.Permission `json:"permissions,omitempty" gorm:"many2many:role_has_permissions;"`
	CreatedAt   *time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
}

type Payload struct {
	UserID uuid.UUID     `json:"user_id"`
	ShopID *uuid.UUID    `json:"shop_id"`
	Roles  []entity.Role `json:"roles"`
}

type CustomClaims struct {
	UserID uuid.UUID     `json:"userId"`
	ShopID *uuid.UUID    `json:"shopId"`
	Roles  []entity.Role `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateAccessToken Generate an access token
func GenerateAccessToken(payload *Payload) (string, error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(commonapp.AccessTokenExpireDuration))

	claims := CustomClaims{
		UserID: payload.UserID,
		ShopID: payload.ShopID,
		Roles:  payload.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

// Generate a refresh token
func GenerateRefreshToken(payload *Payload) (string, error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(commonapp.RefreshTokenExpireDuration))

	claims := CustomClaims{
		UserID: payload.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func VerifyRefreshToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, request.ResponseMessageError{
				Message: "Verify refresh token fail",
			}
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, request.ResponseMessageError{
		Message: "Verify refresh token fail",
	}
}

// GeneratePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// GeneratePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

func GenerateDoubleKey(bitSize int) (*Keys, error) {
	// generate privateKey
	privateKey, errPr := GeneratePrivateKey(bitSize)
	if errPr != nil {
		return nil, errPr
	}

	// generate pubKey
	pubKey, errPub := GeneratePublicKey(&privateKey.PublicKey)
	if errPub != nil {
		return nil, errPub
	}

	return &Keys{
		PublicKey:  string(pubKey),
		PrivateKey: EncodePrivateKeyToPem(privateKey),
	}, nil
}

func EncodePrivateKeyToPem(privateKey *rsa.PrivateKey) string {
	priByte := x509.MarshalPKCS1PrivateKey(privateKey)
	priPem := pem.Block{
		Bytes:   priByte,
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
	}
	return string(pem.EncodeToMemory(&priPem))
}
