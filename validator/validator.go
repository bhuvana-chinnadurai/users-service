package validator

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var (
	isValidFirstName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
	ListOfCountries  = []string{"us", "uk"}
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(strings.Trim(value, " "))
	if n < minLength || n > maxLength {
		return fmt.Errorf("invalid string length, min and max allowed length: %d: %d", minLength, maxLength)
	}
	return nil
}

func ValidateFirstName(firstname string) error {
	if err := ValidateString(firstname, 3, 100); err != nil {
		return err
	}
	if !isValidFirstName(firstname) {
		return fmt.Errorf("provided first_name must contain only letters or spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("provided email address is invalid")
	}
	return nil
}

func ValidateCountry(country string) error {
	if err := ValidateString(country, 2, 2); err != nil {
		return err
	}

	for _, v := range ListOfCountries {
		if v == country {
			return nil
		}
	}

	return fmt.Errorf("provided country is not supported yet")
}
