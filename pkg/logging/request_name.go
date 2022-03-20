package logging

import (
	petname "github.com/dustinkirkland/golang-petname"
)

func UniqueRequestName() string {
	return petname.Generate(2, "-")
}
