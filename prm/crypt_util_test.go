package prm

import (
	"testing"
)

// Test the encoding and decoding of the puffer / uffer

func TestUffer(t *testing.T) {
	uffer := createUffer("test", "0123456789ABCDEF")

	if uffer == "test" {
		t.Error("puffer is equal to 'test', should be encrypted")
	}

	puffer := decryptUffer(uffer, "0123456789ABCDEF")

	if puffer != "test" {
		t.Error("puffer not equal to 'test', got:", puffer)
	}

}

// Test the hash is correctly made for a linux box
// This is a tad hard to test so for now, I just test string length and formatting
//func TestLinuxPassword(t *testing.T) {
//	hash, salt := CreatePasswordHash("test")
//}
