package server

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/oklog/ulid"
	"golang.org/x/crypto/bcrypt"
)

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUuid() (string, error) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		return "", err
	}
	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7f
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xf) | (0x4 << 4)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return uuid, nil
}

func encryptPassword(plain string) (string, error) {
	// 2^10 の強度
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPasswordHash(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err != bcrypt.ErrMismatchedHashAndPassword
}

func createUlid(t time.Time) (string, error) {
	entropy := ulid.Monotonic(rand.Reader, 0)
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
