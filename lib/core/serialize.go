package core

func IndexOf(arr []string, candidate string) int {
	for index, c := range arr {
		if c == candidate {
			return index
		}
	}
	return -1
}
