package auth

import "testing"

func TestSuccess(t *testing.T) {
	token, err := GenerateToken(42, "testuser")
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned an error for a valid token: %v", err)
	}

	if claims.UserId != 42 {
		t.Errorf("Expected UserId 42, got %d", claims.UserId)
	}

	if claims.Username != "testuser" {
		t.Errorf("expected Username 'testuser', got %s", claims.Username)
	}
}

func TestInvalidToken(t *testing.T) {
	_, err := ValidateToken("this is not a real token")
	if err == nil {
		t.Errorf("Expected ValidateToken to fail for a garbage token, but it succeeded")
	}
}
