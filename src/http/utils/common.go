package utils

// FormatErrors puts errors into list under key "errors" for response
func FormatErrors(errors ...error) map[string][]string {
	response := map[string][]string{
		"errors": []string{},
	}

	for _, err := range errors {
		response["errors"] = append(response["errors"], err.Error())
	}

	return response
}
