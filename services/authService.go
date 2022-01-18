package services

import (
	"encoding/hex"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type AuthService struct {
	privKey []byte
	pubKey  []byte

	config config.AuthConfig
}

func NewAuthService(cfg config.AuthConfig) (*AuthService, error) {
	a := AuthService{
		config: cfg,
	}

	seedBytes, err := hex.DecodeString(cfg.JwtKeySeedHex)
	if err != nil {
		return nil, err
	}

	sk := crypto.NewEdKey(seedBytes)

	a.privKey = sk
	a.pubKey = sk[32:]

	return &a, nil
}

func (a *AuthService) CreateToken(pubkey, sig, msg []byte) (string, string, error) {
	err := crypto.VerifySignature(pubkey, msg, sig)
	if err != nil {
		return "", "", err
	}

	bech32Addr := data.NewAddressFromBytes(pubkey)

	jwt, err := a.newJwt(bech32Addr.AddressAsBech32String())
	if err != nil {
		return "", "", err
	}

	refresh, err := crypto.SignPayload(a.privKey, []byte(jwt))

	if err != nil {
		return "", "", err
	}

	return jwt, hex.EncodeToString(refresh), nil
}

func (a *AuthService) RefreshToken(token, refresh string) (string, string, error) {
	claims, err := crypto.GetClaims(token, a.config.JwtSecret, false)
	if err != nil {
		return "", "", err
	}

	refreshBytes, err := hex.DecodeString(refresh)
	if err != nil {
		return "", "", err
	}

	err = crypto.VerifySignature(a.pubKey, []byte(token), refreshBytes)
	if err != nil {
		return "", "", err
	}

	newJwt, err := a.newJwt(claims.Address)
	if err != nil {
		return "", "", err
	}

	newRefresh, err := crypto.SignPayload(a.privKey, []byte(newJwt))
	if err != nil {
		return "", "", err
	}

	return newJwt, hex.EncodeToString(newRefresh), nil
}

func (a *AuthService) newJwt(address string) (string, error) {
	return crypto.GenerateJwt(
		address,
		a.config.JwtSecret,
		a.config.JwtIssuer,
		a.config.JwtExpiryMins,
	)
}
