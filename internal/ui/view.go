package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/codehakase/kelly/internal/ui/components"
	"github.com/codehakase/kelly/pkg/types"
)

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}
	if m.showHelp {
		return components.HelpOverlay(m.width, m.height)
	}

	var sections []string
	sections = append(sections, m.renderTitle(), "", m.renderInputPanel(), "")

	if m.result != nil {
		sections = append(sections, m.renderAllocationBreakdown(), "", m.renderSummary(), "")
	}
	if m.err != nil {
		sections = append(sections, m.renderError(), "")
	}
	sections = append(sections, components.Help())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderTitle() string {
	title := "KELLY • Stake Calculator"
	method := fmt.Sprintf("Method: %s", strings.ToUpper(string(m.method)))

	titleStyle := lipgloss.NewStyle().Foreground(ColorPrimaryText).Bold(true)
	width := m.width
	if width < 80 {
		width = 80
	}

	leftPart := titleStyle.Render(title)
	rightPart := StyleMethod.Render(method)
	spacing := width - lipgloss.Width(leftPart) - lipgloss.Width(rightPart) - 4
	if spacing < 1 {
		spacing = 1
	}

	return lipgloss.NewStyle().
		Background(ColorPanelBG).Padding(0, 2).Width(width).
		Render(leftPart + strings.Repeat(" ", spacing) + rightPart)
}

func (m Model) renderInputPanel() string {
	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Foreground(ColorAccentFocus).Bold(true).Render("INPUT PARAMETERS"))
	sb.WriteString("\n\n")

	colWidth := 35
	leftCol := lipgloss.NewStyle().Width(colWidth)
	rightCol := lipgloss.NewStyle().Width(colWidth)

	leftContent := m.oddsAInput.View() + "\n" + m.nameAInput.View()
	rightContent := m.oddsBInput.View() + "\n" + m.nameBInput.View()
	if m.method == types.MethodKelly {
		leftContent += "\n" + m.probAInput.View()
		rightContent += "\n" + m.probBInput.View()
	}

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftCol.Render(leftContent), "  ", rightCol.Render(rightContent)))
	sb.WriteString("\n\n")
	sb.WriteString(m.totalInput.View())

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).BorderForeground(ColorBorder).Padding(1, 2).
		Render(sb.String())
}

func (m Model) renderAllocationBreakdown() string {
	if m.result == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Foreground(ColorAccentFocus).Bold(true).Render("ALLOCATION BREAKDOWN"))
	sb.WriteString("\n\n")

	colWidth := 35
	leftCol := lipgloss.NewStyle().Width(colWidth)
	rightCol := lipgloss.NewStyle().Width(colWidth)

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		leftCol.Render(m.renderOptionDetails(m.result.OptionA, "OPTION A")), "  ",
		rightCol.Render(m.renderOptionDetails(m.result.OptionB, "OPTION B"))))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).BorderForeground(ColorBorder).Padding(1, 2).
		Render(sb.String())
}

func (m Model) renderOptionDetails(opt types.Option, header string) string {
	var sb strings.Builder

	headerStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(ColorSecondaryText).Width(14)
	valueStyle := lipgloss.NewStyle().Foreground(ColorPrimaryText)

	sb.WriteString(headerStyle.Render(header))
	if opt.Name != "" && opt.Name != header {
		sb.WriteString(" • " + valueStyle.Render(truncateName(opt.Name, 20)))
	}
	sb.WriteString("\n\n")

	sb.WriteString(labelStyle.Render("Odds") + valueStyle.Render(fmt.Sprintf("%.2f ", opt.Odds)) +
		lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("(%.2f%%)", opt.ImpliedProbability*100)) + "\n")

	sb.WriteString(labelStyle.Render("Stake") + valueStyle.Render(fmt.Sprintf("%s%.0f ", m.result.Currency, opt.Stake)) +
		lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("(%.2f%%)", (opt.Stake/m.result.TotalStake)*100)) + "\n")

	sb.WriteString(labelStyle.Render("Return") + valueStyle.Render(fmt.Sprintf("%s%.0f", m.result.Currency, opt.ReturnIfWins)) + "\n")
	sb.WriteString(labelStyle.Render("Profit") + StyleProfit.Render(fmt.Sprintf("+%s%.0f", m.result.Currency, opt.ProfitIfWins)) + "\n")
	sb.WriteString(labelStyle.Render("ROI") + StyleProfit.Render(fmt.Sprintf("+%.2f%%", opt.ROI*100)))

	return sb.String()
}

func (m Model) renderSummary() string {
	if m.result == nil {
		return ""
	}

	labelStyle := lipgloss.NewStyle().Foreground(ColorSecondaryText).Width(22)
	valueStyle := lipgloss.NewStyle().Foreground(ColorPrimaryText)

	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Foreground(ColorAccentFocus).Bold(true).Render("SUMMARY"))
	sb.WriteString("\n\n")

	sb.WriteString(labelStyle.Render("Total Invested") + valueStyle.Render(fmt.Sprintf("%s%.0f", m.result.Currency, m.result.TotalStake)) + "\n")

	sb.WriteString(labelStyle.Render("Guaranteed Profit"))
	if m.result.Summary.GuaranteedProfit {
		sb.WriteString(StyleProfit.Render("YES"))
	} else {
		sb.WriteString(StyleLoss.Render("NO"))
	}
	sb.WriteString("\n")

	sb.WriteString(labelStyle.Render("Profit Range") + StyleProfit.Render(fmt.Sprintf("%s%.0f - %s%.0f",
		m.result.Currency, m.result.Summary.MinProfit, m.result.Currency, m.result.Summary.MaxProfit)) + "\n")

	sb.WriteString(labelStyle.Render("ROI Range") + StyleProfit.Render(fmt.Sprintf("%.2f%% - %.2f%%",
		m.result.Summary.MinROI*100, m.result.Summary.MaxROI*100)) + "\n")

	sb.WriteString(labelStyle.Render("Expected Value") + valueStyle.Render(fmt.Sprintf("%s%.0f (%.2f%%)",
		m.result.Currency, m.result.Summary.ExpectedValue, (m.result.Summary.ExpectedValue/m.result.TotalStake)*100)) + "\n")

	effPct := m.result.Summary.MarketEfficiency * 100
	effStyle := valueStyle
	note := ""
	if effPct < 100 {
		effStyle = StyleProfit
		note = " (Arbitrage opportunity)"
	} else {
		effStyle = StyleLoss
		note = " (No arbitrage)"
	}
	sb.WriteString(labelStyle.Render("Market Efficiency") + effStyle.Render(fmt.Sprintf("%.2f%%", effPct)) +
		lipgloss.NewStyle().Foreground(ColorMuted).Render(note))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).BorderForeground(ColorBorder).Padding(1, 2).
		Render(sb.String())
}

func (m Model) renderError() string {
	return lipgloss.NewStyle().Foreground(ColorLoss).Bold(true).Render("✗ Error: " + m.err.Error())
}

func truncateName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	if maxLen <= 3 {
		return name[:maxLen]
	}
	return name[:maxLen-3] + "..."
}
