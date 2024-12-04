package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func ValidateNationalID(id string) bool {
	// Remove any spaces and ensure length is 10 digits
	id = strings.TrimSpace(id)
	if len(id) != 10 || !isDigits(id) {
		return false
	}
	var digits [10]int
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(id[i]))
		digits[i] = digit
	}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	remainder := sum % 11
	controlDigit := digits[9]

	if remainder < 2 {
		return controlDigit == remainder
	}
	return controlDigit == (11 - remainder)
}
func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

