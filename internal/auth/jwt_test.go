package auth

import (
	"testing"
	"time"
	"fmt"

	"github.com/google/uuid"
)


func TestJWTCreateAndValidate(t *testing.T){
	secretKey := "ayewhadup"
	cases := make([]uuid.UUID, 4)
	for i := 0; i < len(cases); i++ {
		cases[i] = uuid.New()
	}


	for i, c := range cases{
		fmt.Printf("Test [%02d]\n",i)

		// create token
		tokenString, err := MakeJWT(c,secretKey,time.Duration(60*time.Second))
		if err != nil {
			t.Errorf("Error in MakeJWT()")
		}

		// validate token
		userID, err := ValidateJWT(tokenString, secretKey)
		if err != nil {
			t.Errorf("Error in ValidateJWT")
		}

		if userID != c {
			t.Errorf("Somehow the UUIDs do not match")
		}

	}


}


func TestJWTExpiration(t *testing.T){
	secretKey := "ayewhadup"
	userID := uuid.New()

	expiresIn := time.Duration(time.Millisecond*2)
	ticker := time.NewTicker(time.Duration(time.Millisecond*20))

	// Create new token
	tokenString, err := MakeJWT(userID, secretKey, expiresIn)
	if err != nil {
		t.Errorf("Error in MakeJWT()")
	}
	
	// Wait ticker seconds
	_ = <-ticker.C

	_, err = ValidateJWT(tokenString, secretKey)
	if err == nil {
		t.Errorf("Error: Was able to validate after expiration")
	}
}

func TestJWTBadSecretKey(t *testing.T){
	secretKey := "ayewhadup"
	diffSecretKey := "thisisntthesecretkey"
	
	userID := uuid.New()

	expiresIn := time.Duration(time.Millisecond*20)
	ticker := time.NewTicker(time.Duration(time.Millisecond*2))

	// Create new token
	tokenString, err := MakeJWT(userID, secretKey, expiresIn)
	if err != nil {
		t.Errorf("Error in MakeJWT()")
	}
	
	// Wait ticker seconds
	_ = <-ticker.C

	_, err = ValidateJWT(tokenString, diffSecretKey)
	if err == nil {
		t.Errorf("Error: Was able to validate after expiration")
	}
}