package util

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func PrettyNum(n interface{}) string {
	return message.NewPrinter(language.English).Sprintf("%d", n)
}
