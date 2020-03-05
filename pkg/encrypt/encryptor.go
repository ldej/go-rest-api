package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type Encryptor struct {
}

func (e *Encryptor) Encrypt(password string) string {
	result, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(result)
}

func (e *Encryptor) IsValid(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
