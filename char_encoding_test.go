package namaste

import (
	"testing"
)

func TestCharEncoding(t *testing.T) {
	encodedChars := map[string]string{
		"^":  "^5E",
		" ":  "^20",
		"\"": "^22",
		"*":  "^2A",
		"/":  "^2F",
		":":  "^3A",
		"<":  "^3C",
		">":  "^3E",
		"?":  "^3F",
		"\\": "^5C",
		"|":  "^7C",
	}
	decodedChars := map[string]string{
		"^20": " ",
		"^22": "\"",
		"^2A": "*",
		"^2F": "/",
		"^3A": ":",
		"^3C": "<",
		"^3E": ">",
		"^3F": "?",
		"^5C": "\\",
		"^7C": "|",
		"^5E": "^",
	}

	for k, expected := range encodedChars {
		results := charEncode(k)
		if expected != results {
			t.Errorf("expected %q, got %q", expected, results)
		}
	}

	for k, expected := range decodedChars {
		results := charDecode(k)
		if expected != results {
			t.Errorf("expected %q, got %q", expected, results)
		}
	}
}
