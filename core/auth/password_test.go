package auth

import "testing"

func TestDifferentHashFromInput(t *testing.T) {
	hashed, err := HashPass("mypassword123")
	if err != nil {
		t.Fatalf("HashPass returned an error: %v", err)
	}

	if hashed == "mypassword123" {
		t.Errorf("Expected hashed password to differ from plaintext")
	}
}

func TestCorrectPassword(t *testing.T) {
	hashed, err := HashPass("correct-password")
	if err != nil {
		t.Fatalf("HashPass returned an error: %v", err)
	}

	err = CheckPass("correct-password", hashed)
	if err != nil {
		t.Errorf("Expected to succeed for the correct password, got error: %v", err)
	}
}

func TestWrongPassword(t *testing.T) {
	hashed, err := HashPass("correct-password")
	if err != nil {
		t.Fatalf("HashPass returned an error: %v", err)
	}

	err = CheckPass("wrong-password", hashed)
	if err == nil {
		t.Errorf("Expected to fail for an incorrect password, but it succeeded")
	}
}
