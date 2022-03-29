package main

import (
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tchorzewski1991/fds/business/sys/auth"
	"os"
	"strings"
	"time"
)

func main() {

	var sub string
	flag.StringVar(&sub, "sub", "user", "Subject of the JWT (the user)")

	var iss string
	flag.StringVar(&iss, "iss", "fds-toolset", "Issuer of the JWT")

	var dur time.Duration
	flag.DurationVar(&dur, "dur", time.Hour, "Duration of the token validity")

	var perm string
	flag.StringVar(&perm, "perm", "", "Permissions of the user")

	flag.Parse()

	if sub == "" {
		fmt.Println("ERR: subject cannot be empty")
		os.Exit(1)
	}

	if iss == "" {
		fmt.Println("ERR: issuer cannot be empty")
		os.Exit(1)
	}

	if dur == 0 {
		fmt.Println("ERR: duration cannot be 0")
		os.Exit(1)
	}

	perms := strings.Split(perm, ",")

	// Generating a token requires defining a set of claims.
	// For now, we only care about defining the subject and the user.
	// This token will expire in one hour.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Issuer:    iss,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(dur)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Permissions: perms,
	}
	tkn, err := auth.GenerateToken(claims)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tkn)
}
