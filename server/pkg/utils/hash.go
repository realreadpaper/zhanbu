package utils

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCostEnv = "ZHANBU_BCRYPT_COST"

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), configuredBcryptCost())
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func configuredBcryptCost() int {
	value := os.Getenv(bcryptCostEnv)
	if value == "" {
		return bcrypt.DefaultCost
	}

	cost, err := strconv.Atoi(value)
	if err != nil || cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return bcrypt.DefaultCost
	}

	return cost
}

// CheckPassword compares a plaintext password with a bcrypt hash.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
