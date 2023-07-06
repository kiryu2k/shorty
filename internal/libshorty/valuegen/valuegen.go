package valuegen

import (
	"crypto/rand"
	"math/big"
)

const (
	valueSize = 5
	symbols   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var symbolsCount = big.NewInt(int64(len(symbols)))

func GenerateValue() (string, error) {
	value := make([]byte, valueSize)
	for i := range value {

		n, err := rand.Int(rand.Reader, symbolsCount)
		if err != nil {
			return "", err
		}
		value[i] = symbols[int(n.Int64())]
	}
	return string(value), nil
}
