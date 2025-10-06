package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateRandomPassword generates a random 8-digit password
func GenerateRandomPassword() (string, error) {
	const digits = "0123456789"
	const length = 8

	password := make([]byte, length)

	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = digits[num.Int64()]
	}

	return string(password), nil
}

// GenerateRandomPasswordWithPrefix generates a random 8-digit password with an optional prefix
func GenerateRandomPasswordWithPrefix(prefix string) (string, error) {
	password, err := GenerateRandomPassword()
	if err != nil {
		return "", err
	}

	if prefix != "" {
		return fmt.Sprintf("%s%s", prefix, password), nil
	}

	return password, nil
}
