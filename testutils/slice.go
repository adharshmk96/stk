package testutils

func Contains(s []string, e string) bool {
	for _, name := range s {
		if name == e {
			return true
		}
	}
	return false
}
