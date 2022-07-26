package BLC

import (
	"bytes"
	"math/big"
)

// base58 Coding implementation
// 1. Generate a base58 encoding radix table
var b58Alphabet = []byte("" +
	"123456789" +
	"abcdefghijkmnopqrstuvwxyz" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ")

// encoding function
func Base58Encode(input []byte) []byte {
	var result []byte // result
	// big.int
	// byte Convert the byte array tobig.int
	x := big.NewInt(0).SetBytes(input)
	// base length of remainder
	base := big.NewInt(int64(len(b58Alphabet)))
	// find remainder and quotient
	// Judgment condition, whether the final result of removal is 0
	zero := big.NewInt(0)
	// Set the remainder, representing the index position of the base58 radix table
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		// The result obtained is a reversed base58 encoded result
		result = append(result, b58Alphabet[mod.Int64()])
	}
	// return result slice
	Reverse(result)
	// Add a prefix of 1 to represent an address
	result = append([]byte{b58Alphabet[0]}, result...)
	return result

}

// Reverse slice function
func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// Decoding function
// input:base58 Coding results
/*
	Jh83
*/
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	// Prefix index
	zeroBytes := 1
	// 1. Remove prefix
	data := input[zeroBytes:]
	for _, b := range data {
		// 2. Finds the index of the specified number / character in input that appears in the cardinality table(mod)
		charIndex := bytes.IndexByte(b58Alphabet, b) // Internal function that returns the index of the first occurrence of a character in a slice
		// 3. remainder*58
		result.Mul(result, big.NewInt(58))
		// 4. Product result+mod(Index)
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	// 5. Convert to byte array
	decoded := result.Bytes()
	return decoded
}
