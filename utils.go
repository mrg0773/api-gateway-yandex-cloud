package main

func inSlice(searchElement string, slice []string) bool {
	for _, el := range slice {
		if el == searchElement {
			return true
		}
	}
	return false
}
