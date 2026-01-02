package formatter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/codehakase/kelly/pkg/types"
)

func FormatTable(result *types.CalculationResult, verbose bool) string {
	var sb strings.Builder

	sb.WriteString("╭─────────────────────────────────────────────────────────╮\n")

	methodName := strings.Title(string(result.Method))
	title := fmt.Sprintf("│ KELLY • %s Allocation", methodName)
	sb.WriteString(title + strings.Repeat(" ", 58-len(title)) + "│\n")

	sb.WriteString("├─────────────────────────────────────────────────────────┤\n")

	sb.WriteString(fmt.Sprintf("│ %-10s │ Odds: %.2f │ Stake: %s%-6.0f │ +%s%-6.0f │\n",
		truncate(result.OptionA.Name, 10), result.OptionA.Odds,
		result.Currency, result.OptionA.Stake,
		result.Currency, result.OptionA.ProfitIfWins))

	sb.WriteString(fmt.Sprintf("│ %-10s │ Odds: %.2f │ Stake: %s%-6.0f │ +%s%-6.0f │\n",
		truncate(result.OptionB.Name, 10), result.OptionB.Odds,
		result.Currency, result.OptionB.Stake,
		result.Currency, result.OptionB.ProfitIfWins))

	sb.WriteString("├─────────────────────────────────────────────────────────┤\n")

	sb.WriteString(fmt.Sprintf("│ Total: %s%-5.0f │ Profit: %s%-5.0f-%s%-5.0f │ ROI: %.0f-%.0f%% │\n",
		result.Currency, result.TotalStake,
		result.Currency, result.Summary.MinProfit,
		result.Currency, result.Summary.MaxProfit,
		result.Summary.MinROI*100, result.Summary.MaxROI*100))

	sb.WriteString("╰─────────────────────────────────────────────────────────╯\n")

	if verbose {
		sb.WriteString("\n")
		sb.WriteString(formatVerbose(result))
	}

	return sb.String()
}

func formatVerbose(result *types.CalculationResult) string {
	var sb strings.Builder

	sb.WriteString("ℹ Method: ")
	switch result.Method {
	case types.MethodArbitrage:
		sb.WriteString("Arbitrage (Guaranteed Profit)\n")
		sb.WriteString("  Ensures profit regardless of outcome.\n")
	case types.MethodKelly:
		sb.WriteString("Kelly Criterion (Growth Optimization)\n")
		sb.WriteString("  Maximizes long-term growth based on probability estimates.\n")
	case types.MethodProportional:
		sb.WriteString("Proportional (Simple Allocation)\n")
		sb.WriteString("  Allocates stakes inversely to odds.\n")
	}

	sb.WriteString("\nℹ Allocation:\n")
	sb.WriteString(fmt.Sprintf("  - %s: %.2f%%\n", result.OptionA.Name, (result.OptionA.Stake/result.TotalStake)*100))
	sb.WriteString(fmt.Sprintf("  - %s: %.2f%%\n", result.OptionB.Name, (result.OptionB.Stake/result.TotalStake)*100))

	sb.WriteString("\n⚠ Risk:\n")
	if result.Summary.GuaranteedProfit {
		sb.WriteString(fmt.Sprintf("  - Guaranteed profit (efficiency: %.2f%%)\n", result.Summary.MarketEfficiency*100))
	} else {
		sb.WriteString(fmt.Sprintf("  - No guaranteed profit (efficiency: %.2f%%)\n", result.Summary.MarketEfficiency*100))
	}

	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func FormatJSON(result *types.CalculationResult) (string, error) {
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FormatCSV(result *types.CalculationResult) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"Option", "Odds", "Implied_Prob", "Stake", "Return", "Profit", "ROI"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	rowA := []string{
		result.OptionA.Name,
		fmt.Sprintf("%.2f", result.OptionA.Odds),
		fmt.Sprintf("%.2f%%", result.OptionA.ImpliedProbability*100),
		fmt.Sprintf("%.0f", result.OptionA.Stake),
		fmt.Sprintf("%.0f", result.OptionA.ReturnIfWins),
		fmt.Sprintf("%.0f", result.OptionA.ProfitIfWins),
		fmt.Sprintf("%.2f%%", result.OptionA.ROI*100),
	}
	if err := writer.Write(rowA); err != nil {
		return "", err
	}

	rowB := []string{
		result.OptionB.Name,
		fmt.Sprintf("%.2f", result.OptionB.Odds),
		fmt.Sprintf("%.2f%%", result.OptionB.ImpliedProbability*100),
		fmt.Sprintf("%.0f", result.OptionB.Stake),
		fmt.Sprintf("%.0f", result.OptionB.ReturnIfWins),
		fmt.Sprintf("%.0f", result.OptionB.ProfitIfWins),
		fmt.Sprintf("%.2f%%", result.OptionB.ROI*100),
	}
	if err := writer.Write(rowB); err != nil {
		return "", err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
