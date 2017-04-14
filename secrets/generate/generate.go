package generate

import "github.com/untoldwind/trustless/api"

func Password(parameter api.GenerateParameter) (string, error) {
	if parameter.Words != nil {
		return generateWords(parameter.Words)
	}
	if parameter.Chars != nil {
		return generateChars(parameter.Chars)
	}
	return generateChars(defaultCharsParameter)
}
