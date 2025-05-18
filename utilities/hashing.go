package utilities

import "github.com/alexedwards/argon2id"

func ComparePlaintextWithHash(plaintext, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(plaintext, hash)
}

func CreateHash(plaintext string) (string, error) {
	return argon2id.CreateHash(plaintext, argon2id.DefaultParams)
}
