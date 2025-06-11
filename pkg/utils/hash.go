package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPackage interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPwd, plainPwd string) bool
}

type bcryptPackageImpl struct {
}

func NewBcryptPackage() BcryptPackage {
	return &bcryptPackageImpl{}
}

// HashPassword hashes a plain password using bcrypt.
func (b *bcryptPackageImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with a plain password.
func (b *bcryptPackageImpl) ComparePassword(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}
