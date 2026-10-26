package base58

import (
	"encoding/hex"
	"strings"
	"testing"
)

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

func TestEncodeNoSpuriousLeadingOnes(t *testing.T) {
	// Short input that doesn't fill the output buffer — previously
	// produced spurious leading '1's from unused buffer positions
	input := string([]byte{0x41, 0x09, 0xa6})
	encoded := Encode(input)
	if len(encoded) > 0 && encoded[0] == '1' && input[0] != 0x00 {
		t.Errorf("Encode produced spurious leading '1' for non-zero input: '%s'", encoded)
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

func TestDecode(t *testing.T) {
	decoded, err := Decode(TestValue1)
	if err != nil {
		t.Fatalf("Decode(%s) errored: %v", TestValue1, err)
	}
	if decoded != TestString1 {
		t.Errorf("decoded value should be %q. got %q", TestString1, decoded)
	}

	decoded, err = Decode(TestValue2)
	if err != nil {
		t.Fatalf("Decode(%s) errored: %v", TestValue2, err)
	}
	if decoded != TestString2 {
		t.Errorf("decoded value should be %q. got %q", TestString2, decoded)
	}
}

func TestDecodeSolanaKeys(t *testing.T) {
	keys := map[string]string{
		"11111111111111111111111111111111":             "0000000000000000000000000000000000000000000000000000000000000000",
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": "c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d61",
		"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA":  "06ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a9",
		"TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnwqTCSdcq7dc":  "06ddf6e1ee758fde18425dbce46ccddab61afc4d83b90d2905d3b4d3705c9c13",
		"ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL": "8c97258f4e2489f1bb3d1029148e0d830b5a1399daff1084048e7bd8dbe9f859",
	}
	for encoded, expectedHex := range keys {
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Decode(%s) errored: %v", encoded, err)
			continue
		}
		if got := hex.EncodeToString([]byte(decoded)); got != expectedHex {
			t.Errorf("Decode(%s) should be %s. got %s", encoded, expectedHex, got)
		}
		if len(decoded) != 32 {
			t.Errorf("Decode(%s) should be 32 bytes. got %d", encoded, len(decoded))
		}
		if reencoded := EncodeWithLeadingZeros(decoded); reencoded != encoded {
			t.Errorf("EncodeWithLeadingZeros(Decode(%s)) should round-trip. got %s", encoded, reencoded)
		}
	}
}

func TestDecodeRoundTripsLeadingZeros(t *testing.T) {
	inputs := []string{
		string([]byte{0x00, 0x01, 0x02, 0x03}),
		string([]byte{0x00, 0x00, 0x00, 0x01, 0x02}),
		string([]byte{0x00, 0x00, 0x00}),
		string([]byte{0x00, 0xFF, 0xAB, 0xCD}),
		string([]byte{0x41, 0x09, 0xa6}),
		strings.Repeat("\x00", 32),
		strings.Repeat("\xff", 32),
	}
	for _, input := range inputs {
		decoded, err := Decode(EncodeWithLeadingZeros(input))
		if err != nil {
			t.Errorf("Decode errored for %x: %v", input, err)
			continue
		}
		if decoded != input {
			t.Errorf("Decode(EncodeWithLeadingZeros(%x)) should round-trip. got %x", input, decoded)
		}
	}
}

func TestDecodeEmpty(t *testing.T) {
	decoded, err := Decode("")
	if err != nil {
		t.Fatalf("Decode of an empty string errored: %v", err)
	}
	if decoded != "" {
		t.Errorf("Decode of an empty string should be empty. got %q", decoded)
	}
}

func TestDecodeRejectsCharactersOutsideTheAlphabet(t *testing.T) {
	for _, input := range []string{"0", "4Pa7379MLO", "4Pa7379MLI", "l", "4Pa 7379MLc"} {
		if decoded, err := Decode(input); err == nil {
			t.Errorf("Decode(%q) should error. got %q", input, decoded)
		}
	}
}
