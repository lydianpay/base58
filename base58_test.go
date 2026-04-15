package base58

import "testing"

const TestString1 = "Yo Cuz!"
const TestValue1 = "4Pa7379MLc"
const TestString2 = "bf764852-caf7-457c-85b2-b0210f0b0ec2 b85ecfb3-1eae-4b0a-94ea-eb823dd965c8"
const TestValue2 = "7QpECrvkAkFLNrKun6curM4gBndrwyjEZZqh8LbETw9YA4SvXNekBbryrjDsFTe1a4s6NjaWziYcmaDsoLzgy6VtNyLNG9HcKRFm"

func TestAlphabetString(t *testing.T) {
	// This check is to ensure no one accidentally adds/removes characters to the alphabet
	if len(alphabet) != 58 {
		t.Errorf("alphabet length should be 58. got %d", len(alphabet))
	}

	characterCounts := make(map[int32]int)

	for _, char := range alphabet {
		characterCounts[char]++

		if characterCounts[char] > 1 {
			t.Errorf("alphabet contains duplicate characters. too many: %s", string(char))
		}
	}
}

func TestEncode(t *testing.T) {
	encoded := Encode(TestString1)
	if encoded != TestValue1 {
		t.Errorf("encoded value should be %s. got %s", TestValue1, encoded)
	}

	encoded = Encode(TestString2)
	if encoded != TestValue2 {
		t.Errorf("encoded value should be %s. got %s", TestValue2, encoded)
	}
}

func TestEncodeWithLeadingZeros(t *testing.T) {
	// Input with no leading zeros — should match Encode
	input := TestString1
	encoded := EncodeWithLeadingZeros(input)
	expected := Encode(input)
	if encoded != expected {
		t.Errorf("Expected '%s' (same as Encode), got '%s'", expected, encoded)
	}
}

func TestEncodeWithLeadingZerosPreservesZeros(t *testing.T) {
	// Input with one leading zero byte
	input := string([]byte{0x00, 0x01, 0x02, 0x03})
	encoded := EncodeWithLeadingZeros(input)
	if encoded == "" {
		t.Fatalf("Expected non-empty encoded string")
	}
	if encoded[0] != '1' {
		t.Errorf("Expected leading '1' for leading zero byte, got '%c'", encoded[0])
	}
}

func TestEncodeWithLeadingZerosMultipleZeros(t *testing.T) {
	// Input with three leading zero bytes
	input := string([]byte{0x00, 0x00, 0x00, 0x01, 0x02})
	encoded := EncodeWithLeadingZeros(input)
	if len(encoded) < 3 {
		t.Fatalf("Expected at least 3 characters, got %d", len(encoded))
	}
	for i := 0; i < 3; i++ {
		if encoded[i] != '1' {
			t.Errorf("Expected '1' at position %d for leading zero byte, got '%c'", i, encoded[i])
		}
	}
}

func TestEncodeWithLeadingZerosAllZeros(t *testing.T) {
	// Input of all zero bytes — each zero byte adds a leading '1',
	// plus Encode produces its own output for the zero value
	input := string([]byte{0x00, 0x00, 0x00})
	encoded := EncodeWithLeadingZeros(input)
	if len(encoded) < 3 {
		t.Errorf("Expected at least 3 leading '1' characters for three zero bytes, got '%s'", encoded)
	}
	for i := 0; i < 3; i++ {
		if encoded[i] != '1' {
			t.Errorf("Expected '1' at position %d, got '%c'", i, encoded[i])
		}
	}
}

func TestEncodeWithLeadingZerosEmpty(t *testing.T) {
	// Empty input — Encode("") produces a base output, no leading zeros to add
	encoded := EncodeWithLeadingZeros("")
	encodedBase := Encode("")
	if encoded != encodedBase {
		t.Errorf("Expected same as Encode for empty input, got '%s' vs '%s'", encoded, encodedBase)
	}
}

func TestEncodeWithLeadingZerosDeterministic(t *testing.T) {
	input := string([]byte{0x00, 0xFF, 0xAB, 0xCD})
	encoded1 := EncodeWithLeadingZeros(input)
	encoded2 := EncodeWithLeadingZeros(input)
	if encoded1 != encoded2 {
		t.Errorf("Encoding should be deterministic, got '%s' and '%s'", encoded1, encoded2)
	}
}

func TestEncodeWithLeadingZerosDiffersFromEncode(t *testing.T) {
	// For input with leading zeros, EncodeWithLeadingZeros should produce
	// a longer string than Encode since it preserves the zeros
	input := string([]byte{0x00, 0x00, 0x01, 0x02, 0x03})
	withZeros := EncodeWithLeadingZeros(input)
	without := Encode(input)
	if len(withZeros) <= len(without) {
		t.Errorf("EncodeWithLeadingZeros should be longer than Encode for input with leading zeros. Got '%s' vs '%s'", withZeros, without)
	}
}

func TestEncodeRawBinary(t *testing.T) {
	// Raw binary data with bytes > 127 that are invalid UTF-8
	input := string([]byte{0x00, 0x01, 0xFF, 0xC0, 0x80, 0xDE, 0xAD})
	encoded := Encode(input)
	if encoded == "" {
		t.Error("encoding raw binary data should not return an empty string")
	}

	// Encoding the same input twice should produce the same result
	encoded2 := Encode(input)
	if encoded != encoded2 {
		t.Errorf("encoding should be deterministic. got %s and %s", encoded, encoded2)
	}
}
