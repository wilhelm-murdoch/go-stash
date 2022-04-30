package utils

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

func EstimateReadingTime(text string) string {
	pattern, _ := regexp.Compile(`[^a-zA-Z0-9\s]+`)
	words := strings.Fields(pattern.ReplaceAllString(text, ""))

	return fmt.Sprintf("%.0f mins", math.Ceil(float64(len(words))/float64(238)))
}
