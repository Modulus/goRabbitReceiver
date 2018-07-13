package main

import (
	"testing"
)

func TestCreateHashCheckLength(t *testing.T) {
	message := "Vaffler og kaffi"

	hash := CreateHash(message)

	if len(hash) < 32 {
		t.Fatalf("Hash Length is less than 32, is %d\n", len(hash))
	}
}

func TestCreateHashShouldBeEqualForSameString(t *testing.T) {
	message1 := "Testing 12345"
	//message2 := "Testing 12345"

	hash1 := CreateHash(message1)
	hash2 := CreateHash(message1)

	if hash1 == "" {
		t.Fatalf("Failed to create hash1")
	}

	if hash2 == "" {
		t.Fatalf("Failed to create hash2")
	}

	if hash1 != hash2 {
		t.Fatalf("Hashes are not equal!!!")
	}

}

func TestCreateHashShouldBeDifferentForDifferentStrings(t *testing.T) {
	message1 := "Testing 12345"
	message2 := "Testing 54321"

	hash1 := CreateHash(message1)
	hash2 := CreateHash(message2)

	if hash1 == hash2 {
		t.Fatalf("Hashes are equal!!!")
	}
}
