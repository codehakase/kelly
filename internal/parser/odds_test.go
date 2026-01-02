package parser

import (
	"math"
	"testing"
)

func TestParseOdds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		// Decimal format
		{"decimal 2.5", "2.5", 2.5, false},
		{"decimal 3.85", "3.85", 3.85, false},
		{"decimal 2.0", "2.0", 2.0, false},
		{"decimal 1.01", "1.01", 1.01, false},
		{"decimal 10.0", "10.0", 10.0, false},
		{"decimal with spaces", "  2.5  ", 2.5, false},

		// Percentage format
		{"percentage 39%", "39%", 2.564102564102564, false},
		{"percentage 26%", "26%", 3.8461538461538463, false},
		{"percentage 50%", "50%", 2.0, false},
		{"percentage 100%", "100%", 1.0, false},
		{"percentage 10%", "10%", 10.0, false},
		{"percentage with spaces", "  39%  ", 2.564102564102564, false},

		// Fractional format
		{"fractional 3/2", "3/2", 2.5, false},
		{"fractional 5/2", "5/2", 3.5, false},
		{"fractional 1/1", "1/1", 2.0, false},
		{"fractional 2/1", "2/1", 3.0, false},
		{"fractional 1/2", "1/2", 1.5, false},
		{"fractional 10/3", "10/3", 4.333333333333333, false},
		{"fractional with spaces", " 3 / 2 ", 2.5, false},

		// American format
		{"american +250", "+250", 3.5, false},
		{"american -150", "-150", 1.6666666666666667, false},
		{"american +100", "+100", 2.0, false},
		{"american -100", "-100", 2.0, false},
		{"american +200", "+200", 3.0, false},
		{"american -200", "-200", 1.5, false},

		// Error cases
		{"empty string", "", 0, true},
		{"invalid decimal", "abc", 0, true},
		{"invalid percentage", "abc%", 0, true},
		{"zero percentage", "0%", 0, true},
		{"negative percentage", "-10%", 0, true},
		{"over 100 percentage", "150%", 0, true},
		{"invalid fractional", "3/", 0, true},
		{"fractional division by zero", "3/0", 0, true},
		{"invalid american", "+abc", 0, true},
		{"american zero", "+0", 0, true},
		{"decimal less than 1", "0.5", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseOdds(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseOdds(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseOdds(%q) unexpected error: %v", tt.input, err)
				return
			}

			if !floatEquals(result, tt.expected, 0.0001) {
				t.Errorf("ParseOdds(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDecimal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{"valid 2.5", "2.5", 2.5, false},
		{"valid 1.01", "1.01", 1.01, false},
		{"valid 10.0", "10.0", 10.0, false},
		{"invalid text", "abc", 0, true},
		{"less than 1", "0.9", 0, true},
		{"negative", "-2.5", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDecimal(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseDecimal(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parseDecimal(%q) unexpected error: %v", tt.input, err)
				return
			}

			if !floatEquals(result, tt.expected, 0.0001) {
				t.Errorf("parseDecimal(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParsePercentage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{"39%", "39%", 2.564102564102564, false},
		{"26%", "26%", 3.8461538461538463, false},
		{"50%", "50%", 2.0, false},
		{"100%", "100%", 1.0, false},
		{"invalid", "abc%", 0, true},
		{"zero", "0%", 0, true},
		{"negative", "-10%", 0, true},
		{"over 100", "150%", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePercentage(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parsePercentage(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parsePercentage(%q) unexpected error: %v", tt.input, err)
				return
			}

			if !floatEquals(result, tt.expected, 0.0001) {
				t.Errorf("parsePercentage(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseFractional(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{"3/2", "3/2", 2.5, false},
		{"5/2", "5/2", 3.5, false},
		{"1/1", "1/1", 2.0, false},
		{"2/1", "2/1", 3.0, false},
		{"10/3", "10/3", 4.333333333333333, false},
		{"invalid format", "3", 0, true},
		{"division by zero", "3/0", 0, true},
		{"invalid numerator", "abc/2", 0, true},
		{"invalid denominator", "3/abc", 0, true},
		{"negative numerator", "-3/2", 0, true},
		{"negative denominator", "3/-2", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFractional(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFractional(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parseFractional(%q) unexpected error: %v", tt.input, err)
				return
			}

			if !floatEquals(result, tt.expected, 0.0001) {
				t.Errorf("parseFractional(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAmerican(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{"+250", "+250", 3.5, false},
		{"-150", "-150", 1.6666666666666667, false},
		{"+100", "+100", 2.0, false},
		{"-100", "-100", 2.0, false},
		{"+200", "+200", 3.0, false},
		{"-200", "-200", 1.5, false},
		{"invalid", "+abc", 0, true},
		{"zero", "+0", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAmerican(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseAmerican(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parseAmerican(%q) unexpected error: %v", tt.input, err)
				return
			}

			if !floatEquals(result, tt.expected, 0.0001) {
				t.Errorf("parseAmerican(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestImpliedProbability(t *testing.T) {
	tests := []struct {
		name     string
		odds     float64
		expected float64
	}{
		{"odds 2.5", 2.5, 0.4},
		{"odds 3.85", 3.85, 0.25974025974025977},
		{"odds 2.0", 2.0, 0.5},
		{"odds 1.0", 1.0, 1.0},
		{"odds 4.0", 4.0, 0.25},
		{"zero odds", 0.0, 0.0},
		{"negative odds", -2.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ImpliedProbability(tt.odds)

			if !floatEquals(result, tt.expected, 0.0001) {
				t.Errorf("ImpliedProbability(%v) = %v, want %v", tt.odds, result, tt.expected)
			}
		})
	}
}

// floatEquals checks if two floats are equal within a tolerance.
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
