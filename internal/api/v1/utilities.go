package v1

import "fmt"

func checkMandatoryFields(fields ...string) error {
	for _, field := range fields {
		if field == "" {
			return fmt.Errorf("Mandatory field is empty")
		}
	}
	return nil
}
