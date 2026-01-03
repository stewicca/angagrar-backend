package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// ParseSalary parses various salary formats to float64
// Supports: "5 juta", "5jt", "5000000", "Rp 5.000.000", "5.5 juta", etc.
func ParseSalary(input string) (float64, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	// Remove "rp", "rupiah", and spaces
	input = strings.ReplaceAll(input, "rp", "")
	input = strings.ReplaceAll(input, "rupiah", "")
	input = strings.ReplaceAll(input, ".", "")
	input = strings.ReplaceAll(input, ",", ".")
	input = strings.TrimSpace(input)

	// Handle "juta" (millions)
	if strings.Contains(input, "juta") || strings.Contains(input, "jt") {
		input = strings.ReplaceAll(input, "juta", "")
		input = strings.ReplaceAll(input, "jt", "")
		input = strings.TrimSpace(input)

		// Parse the number
		num, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid salary format")
		}

		return num * 1_000_000, nil
	}

	// Handle "miliar" (billions)
	if strings.Contains(input, "miliar") || strings.Contains(input, "m") {
		input = strings.ReplaceAll(input, "miliar", "")
		input = strings.ReplaceAll(input, "m", "")
		input = strings.TrimSpace(input)

		num, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid salary format")
		}

		return num * 1_000_000_000, nil
	}

	// Handle "ribu" (thousands)
	if strings.Contains(input, "ribu") || strings.Contains(input, "rb") {
		input = strings.ReplaceAll(input, "ribu", "")
		input = strings.ReplaceAll(input, "rb", "")
		input = strings.TrimSpace(input)

		num, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid salary format")
		}

		return num * 1_000, nil
	}

	// Direct number parsing
	num, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid salary format")
	}

	return num, nil
}

// NormalizeLocation normalizes city name input to standard format
// Supports variations like: "jakarta", "Jakarta", "jkt", "sby", etc.
func NormalizeLocation(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))

	// Map of abbreviations to full names
	cityMap := map[string]string{
		"jkt":        "Jakarta",
		"jakarta":    "Jakarta",
		"sby":        "Surabaya",
		"surabaya":   "Surabaya",
		"bdg":        "Bandung",
		"bandung":    "Bandung",
		"jogja":      "Yogyakarta",
		"yogya":      "Yogyakarta",
		"yogyakarta": "Yogyakarta",
		"mdn":        "Medan",
		"medan":      "Medan",
		"bali":       "Bali",
		"denpasar":   "Bali",
	}

	if normalized, ok := cityMap[input]; ok {
		return normalized
	}

	// Return capitalized version if not in map
	return strings.Title(input)
}

// NormalizeLifestyle normalizes lifestyle input to standard options
// Returns one of: "Minimalis", "Moderate", "Santai"
func NormalizeLifestyle(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))

	// Minimalis keywords
	minimalisKeywords := []string{
		"minimalis", "hemat", "irit", "sederhana", "saving", "nabung", "1",
	}

	for _, keyword := range minimalisKeywords {
		if strings.Contains(input, keyword) {
			return "Minimalis"
		}
	}

	// Santai keywords
	santaiKeywords := []string{
		"santai", "yolo", "enjoy", "boros", "flexing", "fun", "3",
	}

	for _, keyword := range santaiKeywords {
		if strings.Contains(input, keyword) {
			return "Santai"
		}
	}

	// Moderate keywords (default)
	moderateKeywords := []string{
		"moderate", "balanced", "seimbang", "normal", "biasa", "2",
	}

	for _, keyword := range moderateKeywords {
		if strings.Contains(input, keyword) {
			return "Moderate"
		}
	}

	// Default to Moderate if unclear
	return "Moderate"
}

// ValidateSalary checks if salary is within reasonable bounds
func ValidateSalary(salary float64) error {
	if salary <= 0 {
		return fmt.Errorf("salary must be greater than 0")
	}

	if salary < 1_000_000 {
		return fmt.Errorf("salary seems too low (below 1 million)")
	}

	if salary > 1_000_000_000 {
		return fmt.Errorf("salary seems unrealistically high")
	}

	return nil
}

// IsValidLocation checks if location is in the supported list
func IsValidLocation(location string) bool {
	validLocations := []string{"Jakarta", "Surabaya", "Bandung", "Yogyakarta", "Medan", "Bali"}

	normalized := NormalizeLocation(location)
	for _, valid := range validLocations {
		if normalized == valid {
			return true
		}
	}

	return false
}

// IsValidLifestyle checks if lifestyle is valid
func IsValidLifestyle(lifestyle string) bool {
	validLifestyles := []string{"Minimalis", "Moderate", "Santai"}

	normalized := NormalizeLifestyle(lifestyle)
	for _, valid := range validLifestyles {
		if normalized == valid {
			return true
		}
	}

	return false
}

// RoundToNearest rounds amount to nearest interval (e.g., nearest 1000)
func RoundToNearest(amount float64, interval float64) float64 {
	return math.Round(amount/interval) * interval
}

// ExtractConfirmation checks if user input is a confirmation (yes/no)
func ExtractConfirmation(input string) (bool, bool) {
	input = strings.TrimSpace(strings.ToLower(input))

	// Yes keywords
	yesKeywords := []string{"ya", "yes", "iya", "ok", "oke", "siap", "betul", "benar", "lanjut"}
	for _, keyword := range yesKeywords {
		if strings.Contains(input, keyword) {
			return true, true // isConfirmation=true, isYes=true
		}
	}

	// No keywords
	noKeywords := []string{"tidak", "no", "nope", "enggak", "gak", "salah", "ulang"}
	for _, keyword := range noKeywords {
		if strings.Contains(input, keyword) {
			return true, false // isConfirmation=true, isYes=false
		}
	}

	return false, false // not a clear confirmation
}

// ContainsSalaryInfo checks if input contains salary information using regex
func ContainsSalaryInfo(input string) bool {
	// Regex patterns for salary
	patterns := []string{
		`\d+[\s]*juta`,
		`\d+[\s]*jt`,
		`\d+[\s]*ribu`,
		`\d+[\s]*rb`,
		`\d{7,}`, // 7+ digits (e.g., 5000000)
		`rp[\s]*\d+`,
	}

	input = strings.ToLower(input)
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return true
		}
	}

	return false
}
