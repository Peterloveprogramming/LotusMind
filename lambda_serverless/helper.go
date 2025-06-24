package lambdaServerless

import (
	"unicode"

	"golang.org/x/exp/rand"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateUniqueCode() string {
	result := make([]byte, 11)
	for i := 0; i < 11; i++ {
		num := rand.Intn(len(charset))
		result[i] = charset[num]
	}
	return string(result)
}

func getChakraStatus(score float32) string {
	if score >= 80 && score <= 100 {
		return "Overactive"
	} else if score >= 20 && score < 80 {
		return "Open"
	} else if score >= 0 && score < 20 {
		return "Underactive"
	} else if score >= -50 && score < 0 {
		return "Partially Blocked"
	} else {
		return "Severely Blocked"
	}
}

func containsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
