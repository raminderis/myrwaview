package main

import (
	"os"
	"text/template"
)

type User struct {
	Bio string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	user := User{
		Bio: `<script>alert("haha");</script>`,
	}
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
