package types

import "testing"

func TestNewEmail(t *testing.T) {
	validEmails := []string{
		"test@outlook.com",
		"user.name+tag+sorting@example.com",
		"user_name@example.co.uk",
	}

	for _, email := range validEmails {
		e, err := NewEmail(email)
		if err != nil {
			t.Errorf("NewEmail(%q) returned error: %v", email, err)
		}
		if e.String() != email {
			t.Errorf("NewEmail(%q) = %q; want %q", email, e.String(), email)
		}
	}
}

func TestInvalidEmail(t *testing.T) {
	invalidEmails := []string{
		"plainaddress",
		"missing@domain",
		"@missinglocal.org",
		"user@.invalid.com",
		"",
		"user@invalid..com",
	}

	for _, email := range invalidEmails {
		_, err := NewEmail(email)
		if err == nil {
			t.Errorf("NewEmail(%q) expected error but got nil", email)
		}
	}
}
