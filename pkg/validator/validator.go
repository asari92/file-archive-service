package validator

import (
	"errors"
	"regexp"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func ValidateEmails(emails []string) ([]string, error) {
	validEmails := make([]string, 0, len(emails))
	for _, email := range emails {
		if EmailRX.MatchString(email) {
			validEmails = append(validEmails, email)
		}
	}
	if len(validEmails) == 0 {
		return nil, errors.New("no valid emails provided")
	}

	return validEmails, nil
}
