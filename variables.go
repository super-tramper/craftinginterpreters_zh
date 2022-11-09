package lox

var hadError bool

func GetHadError() bool {
	return hadError
}

func SetHadError(b bool) {
	hadError = b
}
