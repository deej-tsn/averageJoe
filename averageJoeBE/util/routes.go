package util

func ErrorMessage(errorReason string) map[string]string {
	return map[string]string{
		"Error": errorReason,
	}
}
