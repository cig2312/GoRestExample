package shared

import (
	"crypto/rsa"
)

//keys struct to hold pivate and public keys
type AppKeys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}
