package handlers

func formatErrors(errors ...error) map[string][]string {
	response := map[string][]string{
		"errors": []string{},
	}

	for _, err := range errors {
		response["errors"] = append(response["errors"], err.Error())
	}

	return response
}
