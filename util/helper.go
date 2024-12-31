package util

import (
	"fmt"
	"strings"
)

func GetPlatformTypeBasedOnSessionType(sessionType string) string {
	// Split the session type by "_"
	parts := strings.Split(sessionType, "_")

	fmt.Printf("%v", parts)

	// Check if there are any parts and return the last one
	if len(parts) > 0 {
		fmt.Println("part is big than 0")
		return parts[len(parts)-1] // Return the last element
	}

	return "" // Return empty if the sessionType is an empty string
}

func GetDateFormat() string {
	return "2006-01-02"
}
