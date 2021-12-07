package v1

import "fmt"

func IsEmpty(field string) bool {
	return field == ""
}

func checkMandatoryFields(fields ...string) error {
	for _, field := range fields {
		if IsEmpty(field) {
			return fmt.Errorf("Mandatory field is empty")
		}
	}
	return nil
}
