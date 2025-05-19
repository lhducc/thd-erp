package utils

import "regexp"

func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^(?:\+84|0)(?:[1-9])[0-9]{8,9}$`)
	return re.MatchString(phone)
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
