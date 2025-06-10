package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"financeAppAPI/internal/config"
	"io"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/openpgp"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashCVV(cvv string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	return string(bytes), err
}

// PGP encrypt/decrypt for card data
func pgpEntity() (*openpgp.Entity, error) {
	key := os.Getenv("PGP_PRIVATE_KEY")
	pass := os.Getenv("PGP_PASSPHRASE")
	if key == "" {
		return nil, errors.New("PGP_PRIVATE_KEY not set")
	}
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(key))
	if err != nil || len(entityList) == 0 {
		return nil, errors.New("invalid PGP key")
	}
	entity := entityList[0]
	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		if err := entity.PrivateKey.Decrypt([]byte(pass)); err != nil {
			return nil, err
		}
	}
	return entity, nil
}

func EncryptPGP(data string) (string, error) {
	entity, err := pgpEntity()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	w, err := openpgp.Encrypt(&buf, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(data))
	if err != nil {
		return "", err
	}
	w.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecryptPGP(encrypted string) (string, error) {
	entity, err := pgpEntity()
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	md, err := openpgp.ReadMessage(bytes.NewReader(ciphertext), openpgp.EntityList{entity}, nil, nil)
	if err != nil {
		return "", err
	}
	plaintext, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// AES fallback for non-card data
func Encrypt(data string) (string, error) {
	secret := config.LoadConfig().JWTSecret
	if len(secret) < 16 {
		return "", errors.New("секретный ключ слишком короткий для шифрования")
	}
	key := []byte(secret[:16])
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext := []byte(data)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encrypted string) (string, error) {
	secret := config.LoadConfig().JWTSecret
	if len(secret) < 16 {
		return "", errors.New("секретный ключ слишком короткий для расшифровки")
	}
	key := []byte(secret[:16])
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

func JWTSecret() string {
	return config.LoadConfig().JWTSecret
}

func GenerateHMAC(data string) (string, error) {
	secret := config.LoadConfig().JWTSecret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func VerifyHMAC(data, mac string) bool {
	expected, err := GenerateHMAC(data)
	if err != nil {
		return false
	}
	return hmac.Equal([]byte(expected), []byte(mac))
}
