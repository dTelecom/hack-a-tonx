package ton

import (
	"crypto/ed25519"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func getCellHash(message []byte) []byte {
	return cell.BeginCell().
		MustStoreSlice(message, uint(len(message))*8).
		EndCell().
		Hash()
}

func Verify(publicKey ed25519.PublicKey, message, sig []byte) bool {
	return ed25519.Verify(publicKey, getCellHash(message), sig)
}

func Sign(privateKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, getCellHash(message))
}
