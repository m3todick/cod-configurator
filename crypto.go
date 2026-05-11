package main

import (
	"fmt"
	"log"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(plain string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt hash error: %v", err)
	}
	return string(h)
}

func checkPassword(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

func genID() string {
	return fmt.Sprintf("%x", rand.Int63())
}
