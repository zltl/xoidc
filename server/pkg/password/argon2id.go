package password

import (
	"github.com/alexedwards/argon2id"
)

// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html

var DefaultParams = &argon2id.Params{
	Memory:      19 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

// argon2id
func ComparePasswordAndHash(password, hash string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

// argon2id
func CreateHash(password string) (hash string, err error) {
	return argon2id.CreateHash(password, DefaultParams)
}
