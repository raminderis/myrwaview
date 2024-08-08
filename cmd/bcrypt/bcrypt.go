package main

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "aaa"
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashPasswd))
	fmt.Println(base64.URLEncoding.EncodeToString(hashPasswd))
	newpassword := "somepassword"
	err = bcrypt.CompareHashAndPassword(hashPasswd, []byte(newpassword))
	if err != nil {
		fmt.Println("passwd doesnt match")
		panic(err)
	}
	fmt.Println("password matches")
}
