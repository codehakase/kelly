package formatter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/codehakase/kelly/pkg/types"
)

// Helper function to create a sample result
func sampleResult() *types.CalculationResult {
	return &types.CalculationResult{
		Method:     types.MethodArbitrage,
		TotalStake: 10000,
		Currency:   "₦",
		OptionA: types.Option{
			Name:               "Davido - With You",
			Odds:               2.56,
			ImpliedProbability: 0.3906,
			Stake:              6463,
			ReturnIfWins:       16545,
			ProfitIfWins:       6545,
			ROI:                0.6545,
		},
		OptionB: types.Option{
			Name:               "Tyla - PUSH 2 START",
			Odds:               3.85,
			ImpliedProbability: 0.2597,
			Stake:              3537,
			ReturnIfWins:       13617,
			ProfitIfWins:       3617,
			ROI:                0.3617,
		},
		Summary: types.Summary{
			GuaranteedProfit: true,
			MinProfit:        3617,
			MaxProfit:        6545,
			ExpectedValue:    5081,
			MinROI:           0.3617,
			MaxROI:           0.6545,
			MarketEfficiency: 0.6503,
		},
	}
}

func TestFormatTable(t *testing.T) {
	result := sampleResult()

	// Test basic table format
	table := FormatTable(result, false)

	// Check that table contains required elements
	if !strings.Contains(table, "KELLY") {
		t.Error("Table should contain 'KELLY' title")
	}

	if !strings.Contains(table, "Arbitrage") {
		t.Error("Table should contain method name")
	}

	// Check for truncated names (table uses 10 char limit)
	if !strings.Contains(table, "Davido") {
		t.Error("Table should contain part of Option A name")
	}

	if !strings.Contains(table, "Tyla") {
		t.Error("Table should contain part of Option B name")
	}

	// Check for borders
	if !strings.Contains(table, "╭") || !strings.Contains(table, "╰") {
		t.Error("Table should contain top and bottom borders")
	}

	if !strings.Contains(table, "│") {
		t.Error("Table should contain vertical borders")
	}

	// Test verbose mode
	tableVerbose := FormatTable(result, true)

	if !strings.Contains(tableVerbose, "Method:") {
		t.Error("Verbose table should contain method explanation")
	}

	if !strings.Contains(tableVerbose, "Allocation") {
		t.Error("Verbose table should contain allocation section")
	}

	if !strings.Contains(tableVerbose, "Risk") {
		t.Error("Verbose table should contain risk section")
	}
}

func TestFormatJSON(t *testing.T) {
	result := sampleResult()

	jsonStr, err := FormatJSON(result)
	if err != nil {
		t.Fatalf("FormatJSON() error: %v", err)
	}

	// Check that it's valid JSON
	var parsed types.CalculationResult
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("FormatJSON() produced invalid JSON: %v", err)
	}

	// Check that parsed values match original
	if parsed.Method != result.Method {
		t.Errorf("Method = %v, want %v", parsed.Method, result.Method)
	}

	if parsed.TotalStake != result.TotalStake {
		t.Errorf("TotalStake = %.2f, want %.2f", parsed.TotalStake, result.TotalStake)
	}

	if parsed.OptionA.Name != result.OptionA.Name {
		t.Errorf("OptionA.Name = %s, want %s", parsed.OptionA.Name, result.OptionA.Name)
	}

	// Check formatting (should have indentation)
	if !strings.Contains(jsonStr, "  ") {
		t.Error("JSON should be indented")
	}
}

func TestFormatCSV(t *testing.T) {
	result := sampleResult()

	csvStr, err := FormatCSV(result)
	if err != nil {
		t.Fatalf("FormatCSV() error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(csvStr), "\n")

	// Should have 3 lines: header + 2 data rows
	if len(lines) != 3 {
		t.Errorf("CSV should have 3 lines, got %d", len(lines))
	}

	// Check header
	if !strings.Contains(lines[0], "Option") {
		t.Error("CSV header should contain 'Option'")
	}

	if !strings.Contains(lines[0], "Odds") {
		t.Error("CSV header should contain 'Odds'")
	}

	if !strings.Contains(lines[0], "Stake") {
		t.Error("CSV header should contain 'Stake'")
	}

	if !strings.Contains(lines[0], "ROI") {
		t.Error("CSV header should contain 'ROI'")
	}

	// Check that data rows contain option names
	csvContent := strings.Join(lines, "\n")
	if !strings.Contains(csvContent, result.OptionA.Name) {
		t.Errorf("CSV should contain Option A name: %s", result.OptionA.Name)
	}

	if !strings.Contains(csvContent, result.OptionB.Name) {
		t.Errorf("CSV should contain Option B name: %s", result.OptionB.Name)
	}
}

func TestFormatTable_DifferentMethods(t *testing.T) {
	methods := []types.CalculationMethod{
		types.MethodArbitrage,
		types.MethodKelly,
		types.MethodProportional,
	}

	for _, method := range methods {
		result := sampleResult()
		result.Method = method

		table := FormatTable(result, false)

		// Check that method name appears in table
		methodName := strings.Title(string(method))
		if !strings.Contains(table, methodName) {
			t.Errorf("Table should contain method name: %s", methodName)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"short string", "Hello", 10, "Hello"},
		{"exact length", "Hello", 5, "Hello"},
		{"needs truncation", "Hello World", 8, "Hello..."},
		{"very short max", "Hello", 3, "Hel"},
		{"empty string", "", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}
