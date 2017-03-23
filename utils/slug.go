package utils

import "strings"

type Slug struct{}

//Very naive slug generator.
//Just replace all spaces with an hyphen
func (s Slug) Generate(val string) string {
	return strings.Replace(val, " ", "-", len(val))
}

func NewSlugGenerator() Slug {
	return Slug{}
}
