package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// Encrypt 使用 bcrypt 加密纯文本.
func Encrypt(source string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(password), err
}

// Compare 比较密文和明文是否相同.
func Compare(hashdPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashdPassword), []byte(password))
}
