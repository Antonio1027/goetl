package utils

import (
	"strings"
	"time"
	"fmt"
	"math"
	"os"
	"github.com/gin-gonic/gin"
)

// create a function to convert this format 2025-08-10T22:10:00Z to YYYY-MM-DD
func NormalizeDate(dateStr string) (string, error) {
	layouts := []string{
		time.RFC3339,                        // "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05",               // "2006-01-02 15:04:05"
		"2006-01-02",                        // "2006-01-02"
		"2006-01-02 15:04:05.000",         // "2006-01-02 15:04:05.000"
		"2006-01-02T15:04:05.000",         // "2006-01-02T15:04:05.000"
	}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, dateStr)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
	}
	return "", err
	return "", fmt.Errorf("unable to parse date %q with supported layouts: %w", dateStr, err)
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

func SanitizeInt(i int) int {
	return i
}

func SanitizeFloat(f float64) float64 {
	return f
}

func RoundFloat(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(f*shift) / shift
}

func Getenv(key string) string {
	return os.Getenv(key)
}

// ParseQueryInt parses a query param as int, returns fallback if not present or invalid
func ParseQueryInt(c *gin.Context, key string, fallback int) int {
       val := c.Query(key)
       if val == "" {
	       return fallback
       }
       var i int
       _, err := fmt.Sscanf(val, "%d", &i)
       if err != nil {
	       return fallback
       }
       return i
}