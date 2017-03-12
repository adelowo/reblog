package utils

import (
	"regexp"
)

//Minimal email regex
const email_regex = ".+@.+\\..+"

func IsEmail(email string) bool {
	return regexp.MustCompile(email_regex).Match([]byte(email))
}
