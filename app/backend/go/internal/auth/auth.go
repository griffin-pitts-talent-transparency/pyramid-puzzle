package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateAndExtractCognitoSub(r *http.Request) *string {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := jwt.MapClaims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
	token, err := parser.ParseWithClaims(tokenString, claims, GetCognitoPublicKey)

	sub, _ := claims["sub"].(string)
	valid := err == nil &&
		token != nil &&
		token.Valid &&
		claims["token_use"] == "id" &&
		sub != ""

	return map[bool]*string{true: &sub}[valid]
}

func ValidateCognitoUser(r *http.Request) bool {
	result := false
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := jwt.MapClaims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	token, err := parser.ParseWithClaims(tokenString, claims, GetCognitoPublicKey)
	if err == nil && token != nil && token.Valid && claims["token_use"] == "id" {
		result = true
	}

	return result
}

func GetCognitoPublicKey(t *jwt.Token) (any, error) {
	var (
		key any
		err error
	)

	kid, _ := t.Header["kid"].(string)

	jwksURL := "https://cognito-idp.us-east-2.amazonaws.com/us-east-2_iclUrXbex/.well-known/jwks.json"

	resp, httpErr := http.Get(jwksURL)
	if httpErr == nil {
		defer resp.Body.Close()

		var jwks struct {
			Keys []struct {
				Kid string `json:"kid"`
				N   string `json:"n"`
				E   string `json:"e"`
			} `json:"keys"`
		}

		decodeErr := json.NewDecoder(resp.Body).Decode(&jwks)
		if decodeErr == nil {
			found := false

			for _, k := range jwks.Keys {
				if k.Kid == kid {
					nb, _ := base64.RawURLEncoding.DecodeString(k.N)
					eb, _ := base64.RawURLEncoding.DecodeString(k.E)

					n := new(big.Int).SetBytes(nb)
					e := 0
					for _, b := range eb {
						e = e<<8 + int(b)
					}

					key = &rsa.PublicKey{N: n, E: e}
					err = nil
					found = true
					break
				}
			}

			if !found {
				err = jwt.ErrTokenUnverifiable
			}
		} else {
			err = decodeErr
		}
	} else {
		err = httpErr
	}

	return key, err
}
