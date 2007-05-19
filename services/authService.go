package services

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/crypto"
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

	sk := crypto.GibKeySir(seedBytes)

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

	jwt, err := crypto.GenerateJwt(
		bech32Addr.AddressAsBech32String(),
		a.config.JwtSecret,
		a.config.JwtIssuer,
		a.config.JwtExpiryMins,
	)

	if err != nil {
		return "", "", err
	}

	refresh, err := crypto.SignPayload(a.privKey, []byte(jwt))

	if err != nil {
		return "", "", err
	}

	return jwt, hex.EncodeToString(refresh), nil
}

func (a *AuthService) RefreshToken() {}
