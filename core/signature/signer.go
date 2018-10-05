package signature

import "github.com/icloudland/starchain/crypto"

type Signer interface {
	PrivKey() []byte
	PubKey() *crypto.PubKey
}

