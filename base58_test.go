package base58

import "testing"

const TestString1 = "Yo Cuz!"
const TestValue1 = "4Pa7379MLc"
const TestString2 = "bf764852-caf7-457c-85b2-b0210f0b0ec2 b85ecfb3-1eae-4b0a-94ea-eb823dd965c8"
const TestValue2 = "7QpECrvkAkFLNrKun6curM4gBndrwyjEZZqh8LbETw9YA4SvXNekBbryrjDsFTe1a4s6NjaWziYcmaDsoLzgy6VtNyLNG9HcKRFm"

func TestAlphabetString(t *testing.T) {
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
