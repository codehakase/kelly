package validator

import (
	"errors"
	"strings"
	"testing"

	"github.com/codehakase/kelly/pkg/types"
)

func TestValidateOdds(t *testing.T) {
	tests := []struct {
		name    string
		odds    float64
		wantErr bool
	}{
		{"valid 2.5", 2.5, false},
		{"valid 1.01", 1.01, false},
		{"valid 10.0", 10.0, false},
		{"invalid 1.0", 1.0, true},
		{"invalid 0.5", 0.5, true},
		{"invalid 0", 0, true},
		{"invalid negative", -1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOdds(tt.odds)

			if tt.wantErr && err == nil {
				t.Errorf("ValidateOdds(%v) expected error, got nil", tt.odds)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("ValidateOdds(%v) unexpected error: %v", tt.odds, err)
			}
		})
	}
}

func TestValidateProbability(t *testing.T) {
	tests := []struct {
		name    string
		prob    float64
		wantErr bool
	}{
		{"valid 0.5", 0.5, false},
		{"valid 0.1", 0.1, false},
		{"valid 0.9", 0.9, false},
		{"valid 0.55", 0.55, false},
		{"invalid 0", 0, true},
		{"invalid 1", 1.0, true},
		{"invalid 1.5", 1.5, true},
		{"invalid negative", -0.1, true},
		{"invalid 0.0001", 0.0001, false}, // Small but valid
		{"invalid 0.9999", 0.9999, false}, // Close to 1 but valid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProbability(tt.prob)

			if tt.wantErr && err == nil {
				t.Errorf("ValidateProbability(%v) expected error, got nil", tt.prob)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("ValidateProbability(%v) unexpected error: %v", tt.prob, err)
			}
		})
	}
}

func TestValidateTotalStake(t *testing.T) {
	tests := []struct {
		name    string
		total   float64
		wantErr bool
	}{
		{"valid 1000", 1000, false},
		{"valid 10000", 10000, false},
		{"valid 0.01", 0.01, false},
		{"valid 1", 1, false},
		{"invalid 0", 0, true},
		{"invalid negative", -100, true},
		{"invalid -0.01", -0.01, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTotalStake(tt.total)

			if tt.wantErr && err == nil {
				t.Errorf("ValidateTotalStake(%v) expected error, got nil", tt.total)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("ValidateTotalStake(%v) unexpected error: %v", tt.total, err)
			}
		})
	}
}

func TestValidateCalculationInput(t *testing.T) {
	tests := []struct {
		name        string
		input       *types.CalculationInput
		wantErr     bool
		errContains string
	}{
		{
			name: "valid arbitrage input",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      2.56,
				OddsB:      3.85,
				TotalStake: 10000,
				NameA:      "Option A",
				NameB:      "Option B",
				Currency:   "₦",
			},
			wantErr: false,
		},
		{
			name: "valid kelly input with probabilities",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.1,
				OddsB:      3.5,
				TotalStake: 1000,
				ProbA:      0.55,
				ProbB:      0.40,
				NameA:      "Team A",
				NameB:      "Team B",
				Currency:   "$",
			},
			wantErr: false,
		},
		{
			name: "valid proportional input",
			input: &types.CalculationInput{
				Method:     types.MethodProportional,
				OddsA:      2.0,
				OddsB:      3.0,
				TotalStake: 500,
				NameA:      "Option A",
				NameB:      "Option B",
				Currency:   "$",
			},
			wantErr: false,
		},
		{
			name: "invalid odds A",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      0.5,
				OddsB:      2.0,
				TotalStake: 1000,
			},
			wantErr:     true,
			errContains: "Option A",
		},
		{
			name: "invalid odds B",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      2.0,
				OddsB:      0.9,
				TotalStake: 1000,
			},
			wantErr:     true,
			errContains: "Option B",
		},
		{
			name: "invalid total stake",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: -100,
			},
			wantErr:     true,
			errContains: "total stake",
		},
		{
			name: "kelly without probabilities",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				ProbA:      0,
				ProbB:      0,
			},
			wantErr:     true,
			errContains: "requires probability",
		},
		{
			name: "kelly with invalid probability A",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				ProbA:      1.5,
				ProbB:      0.4,
			},
			wantErr:     true,
			errContains: "Option A probability",
		},
		{
			name: "kelly with missing probability B",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				ProbA:      0.5,
				ProbB:      0,
			},
			wantErr:     true,
			errContains: "requires probability",
		},
		{
			name: "kelly with probabilities summing > 1",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				ProbA:      0.7,
				ProbB:      0.6,
			},
			wantErr:     true,
			errContains: "probabilities sum",
		},
		{
			name: "arbitrage with no profit opportunity",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
			},
			wantErr:     true,
			errContains: "no guaranteed profit",
		},
		{
			name: "multiple errors",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      0.5,
				OddsB:      0.9,
				TotalStake: -100,
			},
			wantErr:     true,
			errContains: "multiple validation errors",
		},
		{
			name: "invalid method",
			input: &types.CalculationInput{
				Method:     "invalid",
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
			},
			wantErr:     true,
			errContains: "invalid calculation method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCalculationInput(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateCalculationInput() expected error, got nil")
					return
				}

				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateCalculationInput() error = %v, want error containing %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCalculationInput() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		errors   []error
		expected string
	}{
		{
			name:     "no errors",
			errors:   []error{},
			expected: "",
		},
		{
			name:     "single error",
			errors:   []error{errors.New("single error")},
			expected: "single error",
		},
		{
			name: "multiple errors",
			errors: []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
			},
			expected: "multiple validation errors:\n  1. error 1\n  2. error 2\n  3. error 3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := ValidationError{Errors: tt.errors}
			result := ve.Error()

			if result != tt.expected {
				t.Errorf("ValidationError.Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidateCalculationInputStrict(t *testing.T) {
	// For now, strict validation is the same as normal validation
	input := &types.CalculationInput{
		Method:     types.MethodArbitrage,
		OddsA:      2.56,
		OddsB:      3.85,
		TotalStake: 10000,
	}

	err := ValidateCalculationInputStrict(input)
	if err != nil {
		t.Errorf("ValidateCalculationInputStrict() unexpected error: %v", err)
	}

	// Test with invalid input
	invalidInput := &types.CalculationInput{
		Method:     types.MethodArbitrage,
		OddsA:      0.5,
		OddsB:      2.0,
		TotalStake: 1000,
	}

	err = ValidateCalculationInputStrict(invalidInput)
	if err == nil {
		t.Error("ValidateCalculationInputStrict() expected error for invalid input, got nil")
	}
}
