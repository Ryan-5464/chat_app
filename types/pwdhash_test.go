package types

import (
	"errors"
	xerr "server/xerrors"
	"testing"
)

func TestNewPwdHash(t *testing.T) {
	pwd := []byte("ValidPass123!")

	hash, err := NewPwdHash(pwd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == nil {
		t.Fatalf("expected non-nil hash")
	}
	if string(hash) == string(pwd) {
		t.Errorf("hashed password should not equal plain password")
	}

	if err := hash.Compare(pwd); err != nil {
		t.Errorf("Compare failed: %v", err)
	}

	wrongPwd := []byte("WrongPass123!")
	if err := hash.Compare(wrongPwd); err == nil {
		t.Errorf("Compare should fail for wrong password")
	}
}

func TestNewPwdHash_Errors(t *testing.T) {
	tests := []struct {
		name    string
		pwd     []byte
		wantErr error
	}{
		{
			name:    "too short password",
			pwd:     []byte("short"),
			wantErr: xerr.PwdTooShort,
		},
		{
			name:    "missing uppercase",
			pwd:     []byte("lowercase1!@#"),
			wantErr: xerr.NoUpperCaseChar,
		},
		{
			name:    "missing lowercase",
			pwd:     []byte("UPPERCASE1!@#"),
			wantErr: xerr.NoLowerCaseChar,
		},
		{
			name:    "missing digit",
			pwd:     []byte("NoDigitPass!"),
			wantErr: xerr.NoDigitChar,
		},
		{
			name:    "missing symbol",
			pwd:     []byte("lowercase1ABC"),
			wantErr: xerr.NoSymbolChar,
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			_, err := NewPwdHash(tcase.pwd)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			// Check if error wraps the expected sentinel error
			if !errors.Is(err, tcase.wantErr) {
				t.Errorf("expected error to wrap %v, got %v", tcase.wantErr, err)
			}
		})
	}
}
