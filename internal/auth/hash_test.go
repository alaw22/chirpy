package auth

import (
	"testing"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPasswordHash(t *testing.T){
	cases := []struct {
		input string
		stored string
		err error
	} {
		{
			input: "Thisismypassword",
			stored: "thisismypassword",
			err: bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			input: "!benjamin1!",
			stored: "!benjamin1!",
			err: nil,
		},

	}

	for _, c := range cases{
		// first hash stored
		hash, err := HashPassword(c.stored)
		if err != nil {
			t.Errorf("Error in HashPassword for some reason")
		}

		// Then compare the entered password with the stored hash
		err = CheckPasswordHash(hash,c.input)
		if !errors.Is(err,c.err){
			t.Errorf("Error mismatch")
		}
	}
}

func TestHashPassword(t *testing.T) {
	cases := []struct {
		input string
		err error
	} {
		{
			input: "Idontknowhowlongthisisgoingtobebutthiswillbethelongeststringicanthinkofonthefly",
			err: bcrypt.ErrPasswordTooLong,
		},
		{
			input: "hellothere",
			err: nil,
		},

	}

	for _, c := range cases {
		_, err := HashPassword(c.input)
		if !errors.Is(err,c.err){
			t.Errorf("Error mismatch in HashedPassword test")
		}
	}
}