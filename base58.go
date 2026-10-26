package base58

import (
	"fmt"
	"strings"
)

// PrecalculatedMultiplier is the quotient of math.Log(256)/math.Log(58)
const PrecalculatedMultiplier = 1.36565823731

// PrecalculatedDivisor is the quotient of math.Log(58)/math.Log(256), rounded up
const PrecalculatedDivisor = 0.73225

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Encode requires a string as input and returns a base58 encoded string
func Encode(input string) string {

	bytes := []byte(input)
	inputLen := int(float64(len(bytes)) * PrecalculatedMultiplier)

	out := make([]byte, inputLen+1)

	maxPosition := inputLen

	// Loop over each byte of the input
	for _, b := range bytes {
		position := inputLen
		// Starting at the end of the byte array, calculate the updated character
		for bit := int(b); position > maxPosition || bit != 0; position-- {
			bit = bit + 256*int(out[position])
			// Set the remainder
			out[position] = byte(bit % 58)
			bit /= 58
		}
		maxPosition = position
	}

	var sb strings.Builder
	sb.Grow(len(out))
	for _, char := range out[maxPosition+1:] {
		sb.WriteByte(alphabet[char])
	}

	return sb.String()
}

// EncodeWithLeadingZeros encodes a string to Base58, preserving leading zero bytes
// as '1' characters. This is required for protocols like Bitcoin (Base58Check)
// where leading zero bytes in the input are significant.
func EncodeWithLeadingZeros(input string) string {
	encoded := Encode(input)

	for _, b := range []byte(input) {
		if b != 0 {
			break
		}
		encoded = string(alphabet[0]) + encoded
	}

	return encoded
}

func Decode(input string) (string, error) {
	zeros := 0
	for zeros < len(input) && input[zeros] == alphabet[0] {
		zeros++
	}

	outputLen := int(float64(len(input))*PrecalculatedDivisor) + 1
	out := make([]byte, outputLen)

	maxPosition := outputLen

	for index := zeros; index < len(input); index++ {
		digit := strings.IndexByte(alphabet, input[index])
		if digit < 0 {
			return "", fmt.Errorf("invalid base58 character %q at position %d", input[index], index)
		}
		position := outputLen - 1
		for ; position >= maxPosition || digit != 0; position-- {
			digit += 58 * int(out[position])
			out[position] = byte(digit % 256)
			digit /= 256
		}
		maxPosition = position + 1
	}

	decoded := make([]byte, 0, zeros+outputLen-maxPosition)
	for range zeros {
		decoded = append(decoded, 0)
	}
	decoded = append(decoded, out[maxPosition:]...)

	return string(decoded), nil
}
