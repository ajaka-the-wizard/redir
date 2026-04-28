package utils

import (
	"strings"
	"testing"
)

func TestValidatePublicKey(t *testing.T) {
	t.Parallel()

	validKey := genPublicKey(42)

	tests := []struct {
		name   string
		key    string
		wantID int
		wantOK bool
	}{
		{
			name:   "valid key returns product id",
			key:    validKey,
			wantID: 42,
			wantOK: true,
		},
		{
			name:   "missing product id fails",
			key:    "s" + strings.TrimPrefix(validKey, "42"),
			wantID: 0,
			wantOK: false,
		},
		{
			name:   "invalid uuid segment fails",
			key:    "42snot-a-uuid-not-a-uuid",
			wantID: 0,
			wantOK: false,
		},
		{
			name:   "empty key fails",
			key:    "",
			wantID: 0,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotID, gotOK := ValidatePublicKey(tt.key)
			if gotID != tt.wantID {
				t.Fatalf("ValidatePublicKey(%q) id = %d, want %d", tt.key, gotID, tt.wantID)
			}
			if gotOK != tt.wantOK {
				t.Fatalf("ValidatePublicKey(%q) ok = %v, want %v", tt.key, gotOK, tt.wantOK)
			}
		})
	}
}
