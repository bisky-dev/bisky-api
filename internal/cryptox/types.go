package cryptox

import "crypto/cipher"

type Cipher struct {
	aead cipher.AEAD
}
