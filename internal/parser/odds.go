package parser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParseOdds(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, errors.New("odds cannot be empty")
	}

	if strings.HasSuffix(input, "%") {
		return parsePercentage(input)
	} else if strings.Contains(input, "/") {
		return parseFractional(input)
	} else if strings.HasPrefix(input, "+") || strings.HasPrefix(input, "-") {
		return parseAmerican(input)
	}
	return parseDecimal(input)
}

func parseDecimal(input string) (float64, error) {
	odds, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid decimal odds '%s': %w", input, err)
	}
	if odds < 1.0 {
		return 0, fmt.Errorf("decimal odds must be >= 1.0, got: %s", input)
	}
	return odds, nil
}

func parsePercentage(input string) (float64, error) {
	percentStr := strings.TrimSuffix(input, "%")
	percentStr = strings.TrimSpace(percentStr)

	percentage, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid percentage odds '%s': %w", input, err)
	}
	if percentage <= 0 {
		return 0, fmt.Errorf("percentage must be > 0, got: %s", input)
	}
	if percentage > 100 {
		return 0, fmt.Errorf("percentage must be <= 100, got: %s", input)
	}
	return 100.0 / percentage, nil
}

func parseFractional(input string) (float64, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid fractional odds '%s'", input)
	}

	numerator, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numerator in '%s': %w", input, err)
	}

	denominator, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid denominator in '%s': %w", input, err)
	}

	if denominator == 0 {
		return 0, fmt.Errorf("denominator cannot be zero in '%s'", input)
	}
	if numerator < 0 || denominator < 0 {
		return 0, fmt.Errorf("fractional odds must be positive, got: %s", input)
	}
	return (numerator / denominator) + 1.0, nil
}

func parseAmerican(input string) (float64, error) {
	american, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid American odds '%s': %w", input, err)
	}
	if american == 0 {
		return 0, fmt.Errorf("American odds cannot be zero")
	}

	if american > 0 {
		return (american / 100.0) + 1.0, nil
	}
	return (100.0 / math.Abs(american)) + 1.0, nil
}

func ImpliedProbability(decimalOdds float64) float64 {
	if decimalOdds <= 0 {
		return 0
	}
	return 1.0 / decimalOdds
}
