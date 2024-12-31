package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Create a Random SessionType
func RandomSessionType() string {
	sessionTypes := []string{
		"tibetan_singing_bowl_mr", "tummo_breathing_mr",
	}
	n := len(sessionTypes)

	// Generate random index
	return sessionTypes[rand.Intn(n)]
}

// RandomGender returns a random gender
func RandomGender() string {
	genders := []string{"Male", "Female", "Other"}
	return genders[rand.Intn(len(genders))]
}

// RandomCountryCode returns a random country code from a predefined list
func RandomCountryCode() string {
	countryCodes := []string{"AF", "AL"}
	return countryCodes[rand.Intn(len(countryCodes))]
}

// RandomPlatform returns a random platform
func RandomPlatform() string {
	platforms := []string{"mobile", "mr"}
	return platforms[rand.Intn(len(platforms))]
}

// RandomEmail
func RandomEmail() string {
	return RandomString(10) + "@example.com"
}
