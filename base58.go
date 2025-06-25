package base58

// PrecalculatedMultiplier is the quotient of math.Log(256)/math.Log(58)
const PrecalculatedMultiplier = 1.36565823731

// TODO: add a test to ensure this is always 58 bytes
var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Encode requires a string as input and returns a base58 encoded string
func Encode(input string) (output string) {

	inputLen := int(float64(len(input)) * PrecalculatedMultiplier)

	out := make([]byte, inputLen+1)

	maxPosition := inputLen

	// Loop over each character of the string
	for _, character := range input {
		position := inputLen
		// Starting at the end of the byte array, calculate the updated character
		for bit := character; position > maxPosition || bit != 0; position-- {
			bit = bit + 256*int32(out[position])
			// Set the remainder
			out[position] = byte(bit % 58)
			bit /= 58
		}
		maxPosition = position
	}

	for _, char := range out {
		output += string(alphabet[char])
	}

	return output
}
