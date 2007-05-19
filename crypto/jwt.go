package crypto

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	Address string
	jwt.StandardClaims
}

var isExpired = func(claims JwtClaims) bool {
	return claims.ExpiresAt < time.Now().UTC().Unix()
}

func GenerateJwt(address, secret, issuer string, minsToExpiration int) (string, error) {
	claims := JwtClaims{
		Address: address,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Minute * time.Duration(minsToExpiration)).Unix(),
			Issuer:    issuer,
		},
	}
	payload := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	return payload.SignedString([]byte(secret))
}

func ValidateJwt(signedToken, secret string) (JwtClaims, error) {
	claims, err := parseToken(signedToken, secret)
	if err != nil {
		return JwtClaims{}, err
	}

	if isExpired(*claims) {
		return JwtClaims{}, ErrJwtExpired
	}

	return *claims, nil
}

func GetClaims(signedToken, secret string, verify bool) (JwtClaims, error) {
	var claims *JwtClaims
	var err error
	if !verify {
		claims, err = parseTokenUnverified(signedToken, secret)
	} else {
		claims, err = parseToken(signedToken, secret)
	}

	if err != nil {
		return JwtClaims{}, err
	}

	return *claims, nil
}

func parseToken(signedToken, secret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, ErrJwtParse
	}

	return claims, nil
}

func parseTokenUnverified(signedToken, secret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			if validationErr.Errors&(jwt.ValidationErrorExpired) != 0 && token != nil {
				claims, okCast := token.Claims.(*JwtClaims)
				if !okCast {
					return nil, ErrJwtParse
				}

				return claims, nil
			}
		}

		return nil, err
	}

	return token.Claims.(*JwtClaims), nil
}
