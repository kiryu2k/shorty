package valuegen

import (
	"crypto/sha256"
	"math/big"
	"strconv"

	"github.com/mr-tron/base58/base58"
)

func GenerateValue(source string) string {
	hash := generateHash(source)
	value := new(big.Int).SetBytes(hash).Uint64()
	bytes := []byte(strconv.FormatUint(value, 10))
	return base58.Encode(bytes)[:8]
}

func generateHash(data string) []byte {
	hashAlg := sha256.New()
	hashAlg.Write([]byte(data))
	return hashAlg.Sum(nil)
}
