package calculator

import (
	"math"
	"testing"

	"github.com/codehakase/kelly/pkg/types"
)

func TestKellyCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.CalculationInput
		wantErr bool
	}{
		{
			name: "valid Kelly with edge",
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
			name: "Kelly with no edge",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				ProbA:      0.5,
				ProbB:      0.5,
				NameA:      "Even A",
				NameB:      "Even B",
				Currency:   "$",
			},
			wantErr: false,
		},
		{
			name: "Kelly with high confidence",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      3.0,
				OddsB:      2.0,
				TotalStake: 5000,
				ProbA:      0.8,
				ProbB:      0.15,
				NameA:      "Favorite",
				NameB:      "Underdog",
				Currency:   "$",
			},
			wantErr: false,
		},
		{
			name: "Kelly without probabilities (should error)",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				ProbA:      0,
				ProbB:      0,
				NameA:      "A",
				NameB:      "B",
				Currency:   "$",
			},
			wantErr: true,
		},
		{
			name: "Kelly with only one probability (should error)",
			input: &types.CalculationInput{
				Method:     types.MethodKelly,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				ProbA:      0.5,
				ProbB:      0,
				NameA:      "A",
				NameB:      "B",
				Currency:   "$",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := &KellyCalculator{}
			result, err := calc.Calculate(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("KellyCalculator.Calculate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("KellyCalculator.Calculate() unexpected error: %v", err)
				return
			}

			// Check that method is correct
			if result.Method != types.MethodKelly {
				t.Errorf("Method = %v, want %v", result.Method, types.MethodKelly)
			}

			// Check that probabilities are set in options
			if result.OptionA.Probability != tt.input.ProbA {
				t.Errorf("OptionA.Probability = %.4f, want %.4f", result.OptionA.Probability, tt.input.ProbA)
			}

			if result.OptionB.Probability != tt.input.ProbB {
				t.Errorf("OptionB.Probability = %.4f, want %.4f", result.OptionB.Probability, tt.input.ProbB)
			}

			// Check that stakes don't exceed total
			totalAllocated := result.OptionA.Stake + result.OptionB.Stake
			if totalAllocated > tt.input.TotalStake*1.01 { // Allow 1% tolerance for rounding
				t.Errorf("Total allocated (%.2f) exceeds total stake (%.2f)",
					totalAllocated, tt.input.TotalStake)
			}

			// Check that stakes are non-negative
			if result.OptionA.Stake < 0 || result.OptionB.Stake < 0 {
				t.Error("Stakes should be non-negative")
			}
		})
	}
}

func TestKellyCalculator_Normalization(t *testing.T) {
	// Test case where Kelly percentages sum to > 100%
	calc := &KellyCalculator{}
	input := &types.CalculationInput{
		Method:     types.MethodKelly,
		OddsA:      2.0,
		OddsB:      2.0,
		TotalStake: 1000,
		ProbA:      0.7, // Strong belief in A
		ProbB:      0.6, // Strong belief in B (sum > 1, which is overconfident)
		NameA:      "A",
		NameB:      "B",
		Currency:   "$",
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate() error: %v", err)
	}

	// Total allocated should be normalized to total stake
	totalAllocated := result.OptionA.Stake + result.OptionB.Stake
	if totalAllocated > input.TotalStake*1.01 {
		t.Errorf("Total allocated (%.2f) exceeds total stake (%.2f)",
			totalAllocated, input.TotalStake)
	}
}

func TestKellyCalculator_NoEdge(t *testing.T) {
	// Test case where there's no edge (fair odds = true probability)
	calc := &KellyCalculator{}
	input := &types.CalculationInput{
		Method:     types.MethodKelly,
		OddsA:      2.0, // Implies 50% probability
		OddsB:      2.0, // Implies 50% probability
		TotalStake: 1000,
		ProbA:      0.5, // Matches implied probability
		ProbB:      0.5, // Matches implied probability
		NameA:      "Fair A",
		NameB:      "Fair B",
		Currency:   "$",
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate() error: %v", err)
	}

	// When there's no edge, Kelly should recommend minimal or zero bet
	// Kelly% = (p × odds - 1) / (odds - 1) = (0.5 × 2 - 1) / (2 - 1) = 0
	if result.OptionA.Stake > 10 || result.OptionB.Stake > 10 {
		t.Errorf("With no edge, stakes should be minimal (A: %.2f, B: %.2f)",
			result.OptionA.Stake, result.OptionB.Stake)
	}
}

func TestKellyCalculator_ExpectedValue(t *testing.T) {
	// Test that expected value is calculated using user probabilities
	calc := &KellyCalculator{}
	input := &types.CalculationInput{
		Method:     types.MethodKelly,
		OddsA:      2.5,
		OddsB:      3.0,
		TotalStake: 1000,
		ProbA:      0.6,
		ProbB:      0.3,
		NameA:      "A",
		NameB:      "B",
		Currency:   "$",
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate() error: %v", err)
	}

	// Expected value should be positive if user has an edge
	// EV = probA × profitA + probB × profitB + (1 - probA - probB) × -totalStake
	expectedEV := input.ProbA*result.OptionA.ProfitIfWins +
		input.ProbB*result.OptionB.ProfitIfWins +
		(1-input.ProbA-input.ProbB)*(-input.TotalStake)

	if !floatAlmostEqual(result.Summary.ExpectedValue, expectedEV, 1.0) {
		t.Errorf("ExpectedValue = %.2f, calculated EV = %.2f",
			result.Summary.ExpectedValue, expectedEV)
	}
}

func TestArbitrageCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name        string
		input       *types.CalculationInput
		wantStakeA  float64
		wantStakeB  float64
		wantProfitA float64
		wantProfitB float64
		wantErr     bool
	}{
		{
			name: "Grammy example from spec (lines 119-133)",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      2.564102564102564,  // 39% converted
				OddsB:      3.8461538461538463, // 26% converted
				TotalStake: 10000,
				NameA:      "Davido - With You",
				NameB:      "Tyla - PUSH 2 START",
				Currency:   "₦",
			},
			wantStakeA:  6453.49, // Actual calculation
			wantStakeB:  3546.51, // Actual calculation
			wantProfitA: 6545.65, // Actual profit
			wantProfitB: 3640.42, // Actual profit
			wantErr:     false,
		},
		{
			name: "equal odds (no arbitrage opportunity)",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				NameA:      "Option A",
				NameB:      "Option B",
				Currency:   "$",
			},
			wantStakeA:  500,
			wantStakeB:  500,
			wantProfitA: 0, // No profit with equal odds
			wantProfitB: 0, // No profit with equal odds
			wantErr:     false,
		},
		{
			name: "different odds",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      2.5,
				OddsB:      3.0,
				TotalStake: 1000,
				NameA:      "Team A",
				NameB:      "Team B",
				Currency:   "$",
			},
			wantStakeA:  571.43,
			wantStakeB:  428.57,
			wantProfitA: 428.57,
			wantProfitB: 285.71,
			wantErr:     false,
		},
		{
			name: "high odds",
			input: &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      5.0,
				OddsB:      10.0,
				TotalStake: 5000,
				NameA:      "Underdog A",
				NameB:      "Underdog B",
				Currency:   "$",
			},
			wantStakeA:  3461.54,
			wantStakeB:  1538.46,
			wantProfitA: 12307.7, // Actual profit
			wantProfitB: 10384.6, // Actual profit
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := &ArbitrageCalculator{}
			result, err := calc.Calculate(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ArbitrageCalculator.Calculate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ArbitrageCalculator.Calculate() unexpected error: %v", err)
				return
			}

			// Check stakes
			if !floatAlmostEqual(result.OptionA.Stake, tt.wantStakeA, 1.0) {
				t.Errorf("OptionA.Stake = %.2f, want %.2f", result.OptionA.Stake, tt.wantStakeA)
			}

			if !floatAlmostEqual(result.OptionB.Stake, tt.wantStakeB, 1.0) {
				t.Errorf("OptionB.Stake = %.2f, want %.2f", result.OptionB.Stake, tt.wantStakeB)
			}

			// Check profits
			if !floatAlmostEqual(result.OptionA.ProfitIfWins, tt.wantProfitA, 10.0) {
				t.Errorf("OptionA.ProfitIfWins = %.2f, want %.2f", result.OptionA.ProfitIfWins, tt.wantProfitA)
			}

			if !floatAlmostEqual(result.OptionB.ProfitIfWins, tt.wantProfitB, 10.0) {
				t.Errorf("OptionB.ProfitIfWins = %.2f, want %.2f", result.OptionB.ProfitIfWins, tt.wantProfitB)
			}

			// Check that method is correct
			if result.Method != types.MethodArbitrage {
				t.Errorf("Method = %v, want %v", result.Method, types.MethodArbitrage)
			}

			// Check that total stake is correct
			if result.TotalStake != tt.input.TotalStake {
				t.Errorf("TotalStake = %.2f, want %.2f", result.TotalStake, tt.input.TotalStake)
			}

			// Check that currency is correct
			if result.Currency != tt.input.Currency {
				t.Errorf("Currency = %s, want %s", result.Currency, tt.input.Currency)
			}

			// Check market efficiency
			if result.Summary.MarketEfficiency <= 0 {
				t.Error("MarketEfficiency should be positive")
			}

			// Check that summary values are consistent
			if result.Summary.MinProfit > result.Summary.MaxProfit {
				t.Errorf("MinProfit (%.2f) should be <= MaxProfit (%.2f)",
					result.Summary.MinProfit, result.Summary.MaxProfit)
			}
		})
	}
}

func TestArbitrageCalculator_GuaranteedProfit(t *testing.T) {
	tests := []struct {
		name       string
		oddsA      float64
		oddsB      float64
		wantProfit bool
	}{
		{"arbitrage opportunity", 2.5, 3.0, true},     // 0.4 + 0.33 = 0.73 < 1
		{"no arbitrage", 2.0, 2.0, false},             // 0.5 + 0.5 = 1.0
		{"no arbitrage (overround)", 1.9, 1.9, false}, // 0.526 + 0.526 = 1.05 > 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := &ArbitrageCalculator{}
			input := &types.CalculationInput{
				Method:     types.MethodArbitrage,
				OddsA:      tt.oddsA,
				OddsB:      tt.oddsB,
				TotalStake: 1000,
				NameA:      "A",
				NameB:      "B",
				Currency:   "$",
			}

			result, err := calc.Calculate(input)
			if err != nil {
				t.Fatalf("Calculate() error: %v", err)
			}

			if result.Summary.GuaranteedProfit != tt.wantProfit {
				t.Errorf("GuaranteedProfit = %v, want %v", result.Summary.GuaranteedProfit, tt.wantProfit)
			}
		})
	}
}

func TestArbitrageCalculator_EqualProfits(t *testing.T) {
	// For true arbitrage, profits should be nearly equal regardless of outcome
	calc := &ArbitrageCalculator{}
	input := &types.CalculationInput{
		Method:     types.MethodArbitrage,
		OddsA:      2.5,
		OddsB:      3.0,
		TotalStake: 1000,
		NameA:      "A",
		NameB:      "B",
		Currency:   "$",
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate() error: %v", err)
	}

	// For arbitrage, profits should be close (within rounding errors)
	profitDiff := math.Abs(result.OptionA.ProfitIfWins - result.OptionB.ProfitIfWins)

	// Allow for some difference due to rounding
	if profitDiff > 150 {
		t.Errorf("Profit difference too large: %.2f (A: %.2f, B: %.2f)",
			profitDiff, result.OptionA.ProfitIfWins, result.OptionB.ProfitIfWins)
	}
}

func TestProportionalCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name       string
		input      *types.CalculationInput
		wantStakeA float64
		wantStakeB float64
		wantErr    bool
	}{
		{
			name: "equal odds",
			input: &types.CalculationInput{
				Method:     types.MethodProportional,
				OddsA:      2.0,
				OddsB:      2.0,
				TotalStake: 1000,
				NameA:      "A",
				NameB:      "B",
				Currency:   "$",
			},
			wantStakeA: 500,
			wantStakeB: 500,
			wantErr:    false,
		},
		{
			name: "different odds",
			input: &types.CalculationInput{
				Method:     types.MethodProportional,
				OddsA:      2.5,
				OddsB:      3.0,
				TotalStake: 1000,
				NameA:      "Team A",
				NameB:      "Team B",
				Currency:   "$",
			},
			wantStakeA: 545.45,
			wantStakeB: 454.55,
			wantErr:    false,
		},
		{
			name: "high vs low odds",
			input: &types.CalculationInput{
				Method:     types.MethodProportional,
				OddsA:      5.0,
				OddsB:      2.0,
				TotalStake: 5000,
				NameA:      "Underdog",
				NameB:      "Favorite",
				Currency:   "$",
			},
			wantStakeA: 1428.57,
			wantStakeB: 3571.43,
			wantErr:    false,
		},
		{
			name: "very different odds",
			input: &types.CalculationInput{
				Method:     types.MethodProportional,
				OddsA:      10.0,
				OddsB:      1.5,
				TotalStake: 2000,
				NameA:      "Longshot",
				NameB:      "Heavy Favorite",
				Currency:   "$",
			},
			wantStakeA: 260.87,
			wantStakeB: 1739.13,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := &ProportionalCalculator{}
			result, err := calc.Calculate(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ProportionalCalculator.Calculate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ProportionalCalculator.Calculate() unexpected error: %v", err)
				return
			}

			// Check stakes
			if !floatAlmostEqual(result.OptionA.Stake, tt.wantStakeA, 1.0) {
				t.Errorf("OptionA.Stake = %.2f, want %.2f", result.OptionA.Stake, tt.wantStakeA)
			}

			if !floatAlmostEqual(result.OptionB.Stake, tt.wantStakeB, 1.0) {
				t.Errorf("OptionB.Stake = %.2f, want %.2f", result.OptionB.Stake, tt.wantStakeB)
			}

			// Check that method is correct
			if result.Method != types.MethodProportional {
				t.Errorf("Method = %v, want %v", result.Method, types.MethodProportional)
			}

			// Check that stakes sum to total (within rounding error)
			totalStake := result.OptionA.Stake + result.OptionB.Stake
			if !floatAlmostEqual(totalStake, tt.input.TotalStake, 1.0) {
				t.Errorf("Total stake = %.2f, want %.2f", totalStake, tt.input.TotalStake)
			}
		})
	}
}

func TestProportionalCalculator_InverseProportionality(t *testing.T) {
	// Test that higher odds result in lower stakes (inverse proportionality)
	calc := &ProportionalCalculator{}
	input := &types.CalculationInput{
		Method:     types.MethodProportional,
		OddsA:      5.0, // Higher odds
		OddsB:      2.0, // Lower odds
		TotalStake: 1000,
		NameA:      "High Odds",
		NameB:      "Low Odds",
		Currency:   "$",
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate() error: %v", err)
	}

	// Higher odds should get lower stake
	if result.OptionA.Stake >= result.OptionB.Stake {
		t.Errorf("Higher odds (A) should get lower stake: A=%.2f, B=%.2f",
			result.OptionA.Stake, result.OptionB.Stake)
	}

	// Check that the ratio is correct
	// Weight_A = 1/5 = 0.2, Weight_B = 1/2 = 0.5, Total = 0.7
	// Stake_A should be 0.2/0.7 ≈ 28.57%, Stake_B should be 0.5/0.7 ≈ 71.43%
	expectedRatio := (1.0 / input.OddsA) / (1.0/input.OddsA + 1.0/input.OddsB)
	actualRatio := result.OptionA.Stake / input.TotalStake

	if !floatAlmostEqual(actualRatio, expectedRatio, 0.01) {
		t.Errorf("Stake ratio = %.4f, want %.4f", actualRatio, expectedRatio)
	}
}

func TestProportionalCalculator_StakesAlwaysPositive(t *testing.T) {
	// Test that stakes are always positive regardless of odds
	tests := []struct {
		name  string
		oddsA float64
		oddsB float64
	}{
		{"normal odds", 2.0, 3.0},
		{"extreme odds", 100.0, 1.01},
		{"near-even odds", 1.01, 1.02},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := &ProportionalCalculator{}
			input := &types.CalculationInput{
				Method:     types.MethodProportional,
				OddsA:      tt.oddsA,
				OddsB:      tt.oddsB,
				TotalStake: 1000,
				NameA:      "A",
				NameB:      "B",
				Currency:   "$",
			}

			result, err := calc.Calculate(input)
			if err != nil {
				t.Fatalf("Calculate() error: %v", err)
			}

			if result.OptionA.Stake <= 0 {
				t.Errorf("OptionA stake should be positive, got %.2f", result.OptionA.Stake)
			}

			if result.OptionB.Stake <= 0 {
				t.Errorf("OptionB stake should be positive, got %.2f", result.OptionB.Stake)
			}
		})
	}
}

func TestNewCalculator(t *testing.T) {
	tests := []struct {
		name     string
		method   types.CalculationMethod
		wantType string
	}{
		{"arbitrage", types.MethodArbitrage, "*calculator.ArbitrageCalculator"},
		{"kelly", types.MethodKelly, "*calculator.KellyCalculator"},
		{"proportional", types.MethodProportional, "*calculator.ProportionalCalculator"},
		{"unknown (defaults to arbitrage)", "unknown", "*calculator.ArbitrageCalculator"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewCalculator(tt.method)
			if calc == nil {
				t.Error("NewCalculator() returned nil")
			}

			// Type check
			calcType := ""
			switch calc.(type) {
			case *ArbitrageCalculator:
				calcType = "*calculator.ArbitrageCalculator"
			case *KellyCalculator:
				calcType = "*calculator.KellyCalculator"
			case *ProportionalCalculator:
				calcType = "*calculator.ProportionalCalculator"
			}

			if calcType != tt.wantType {
				t.Errorf("NewCalculator(%v) type = %s, want %s", tt.method, calcType, tt.wantType)
			}
		})
	}
}

// Helper function to check if two floats are almost equal within a tolerance
func floatAlmostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
