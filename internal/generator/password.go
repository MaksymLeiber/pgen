package generator

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/MaksymLeiber/pgen/internal/i18n"
	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 16
	// Набор для генерации
	charsetAlphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	charsetSymbols  = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	charsetFull     = charsetAlphaNum + charsetSymbols
)

// сонфиг Argon2
type ArgonConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

type PasswordGenerator struct {
	length int
	argon  *ArgonConfig
}

func NewPasswordGenerator(length int) *PasswordGenerator {
	if length <= 0 {
		length = 16
	}
	return &PasswordGenerator{
		length: length,
		argon:  nil,
	}
}

// NewPasswordGeneratorWithConfig создаёт генератор с пользовательской конфигурацией Argon2
func NewPasswordGeneratorWithConfig(length int, argonConfig ArgonConfig) *PasswordGenerator {
	if length <= 0 {
		length = 16
	}
	return &PasswordGenerator{
		length: length,
		argon:  &argonConfig,
	}
}

func (pg *PasswordGenerator) GeneratePassword(masterPassword, serviceName, username, pepper string, messages *i18n.Messages) (string, error) {
	salt := createSalt(serviceName, username, pepper)

	// Определяем параметрыы аrgon2
	var argonTime uint32 = 3
	var argonMemory uint32 = 256 * 1024
	var argonThreads uint8 = 4
	var argonKeyLen uint32 = 32

	if pg.argon != nil {
		argonTime = pg.argon.Time
		argonMemory = pg.argon.Memory
		argonThreads = pg.argon.Threads
		argonKeyLen = pg.argon.KeyLen
	}

	var keyLen uint32 = argonKeyLen
	if pg.length > int(argonKeyLen) {
		keyLen = uint32(pg.length * 2) // Удваиваем для запаса
	}

	hash := argon2.IDKey(
		[]byte(masterPassword),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		keyLen,
	)

	// Используем все биты хеша
	password := pg.generateFromHash(hash)

	if len(password) < pg.length {
		return "", fmt.Errorf("%s", messages.Errors.HashTooShort)
	}

	return password[:pg.length], nil
}

// generateFromHash генерирует пароль из хеша без потери энтропии
func (pg *PasswordGenerator) generateFromHash(hash []byte) string {
	charset := charsetFull
	charsetLen := big.NewInt(int64(len(charset)))

	hashInt := new(big.Int).SetBytes(hash)

	password := make([]byte, 0, pg.length*2)

	// Генерируем используя все биты хеша
	for hashInt.Sign() > 0 && len(password) < pg.length*2 {
		remainder := new(big.Int)
		hashInt.DivMod(hashInt, charsetLen, remainder)
		password = append(password, charset[remainder.Int64()])
	}

	// Если хеш закончился, а пароль недостаточно длинный, расширяем его >
	for len(password) < pg.length {
		// Используем SHA256 от текущего пароля для расширения
		newHash := sha256.Sum256(append(hash, password...))
		newHashInt := new(big.Int).SetBytes(newHash[:])

		for newHashInt.Sign() > 0 && len(password) < pg.length {
			remainder := new(big.Int)
			newHashInt.DivMod(newHashInt, charsetLen, remainder)
			password = append(password, charset[remainder.Int64()])
		}
		hash = newHash[:]
	}

	return string(password)
}

func createSalt(serviceName, username, pepper string) []byte {
	// Улучшенная генерация salt с версионированием и дополнительной энтропией
	baseText := "PGenCLI|v1|" + serviceName + "|" + username + "|" + pepper
	hash := sha256.Sum256([]byte(baseText))
	return hash[:saltLength]
}

func ValidateLength(length int) error {
	if length < 4 {
		return fmt.Errorf("length_too_short")
	}
	if length > 128 {
		return fmt.Errorf("length_too_long")
	}
	return nil
}
