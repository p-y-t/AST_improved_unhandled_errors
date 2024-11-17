package utils

// ReverseString returns the reversed form of the input string.
func ReverseString(s string) (string, int, error) {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), len(string(runes)), nil
}
