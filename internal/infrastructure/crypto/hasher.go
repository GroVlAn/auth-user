package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/GroVlAn/auth-user/internal/core/e"
	"golang.org/x/crypto/argon2"
)

type Deps struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

type Argon2Hasher struct {
	Deps
}

func New(deps Deps) *Argon2Hasher {
	return &Argon2Hasher{
		Deps: deps,
	}
}

func (a Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, a.SaltLen)
	_, _ = rand.Read(salt)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a.Time,
		a.Memory,
		a.Threads,
		a.KeyLen,
	)

	return base64.StdEncoding.EncodeToString(append(salt, hash...)), nil
}

func (a Argon2Hasher) Compare(encodedHash, password string) error {
	data, _ := base64.StdEncoding.DecodeString(encodedHash)

	salt := data[:a.SaltLen]
	hash := data[a.SaltLen:]

	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		a.Time,
		a.Memory,
		a.Threads,
		a.KeyLen,
	)

	if string(hash) != string(newHash) {
		return e.ErrPasswordMismatch
	}

	return nil
}
