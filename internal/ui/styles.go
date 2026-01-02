package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBackground    = lipgloss.Color("#0a0e27")
	ColorPanelBG       = lipgloss.Color("#1a1e3f")
	ColorBorder        = lipgloss.Color("#2d3561")
	ColorPrimaryText   = lipgloss.Color("#e4e4e7")
	ColorSecondaryText = lipgloss.Color("#9ca3af")
	ColorMuted         = lipgloss.Color("#6b7280")
	ColorAccentFocus   = lipgloss.Color("#60a5fa")
	ColorHighlight     = lipgloss.Color("#f59e0b")
	ColorProfit        = lipgloss.Color("#10b981")
	ColorLoss          = lipgloss.Color("#ef4444")
)

var (
	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorPrimaryText).
			Background(ColorPanelBG).
			Padding(0, 2).
			Bold(true)

	StylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	StylePanelTitle = lipgloss.NewStyle().
			Foreground(ColorAccentFocus).
			Bold(true)

	StyleInputLabel = lipgloss.NewStyle().
			Foreground(ColorSecondaryText).
			Width(12)

	StyleInputActive = lipgloss.NewStyle().
				Foreground(ColorAccentFocus).
				Bold(true)

	StyleInputInactive = lipgloss.NewStyle().
				Foreground(ColorPrimaryText)

	StyleInputPlaceholder = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Italic(true)

	StyleInputError = lipgloss.NewStyle().
			Foreground(ColorLoss)

	StyleProfit = lipgloss.NewStyle().
			Foreground(ColorProfit).
			Bold(true)

	StyleLoss = lipgloss.NewStyle().
			Foreground(ColorLoss).
			Bold(true)

	StyleHighlight = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleHelpKey = lipgloss.NewStyle().
			Foreground(ColorAccentFocus).
			Bold(true)

	StyleMethod = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)

	StyleTableHeader = lipgloss.NewStyle().
				Foreground(ColorSecondaryText).
				Bold(true)

	StyleTableValue = lipgloss.NewStyle().
			Foreground(ColorPrimaryText)

	StyleTableLabel = lipgloss.NewStyle().
			Foreground(ColorSecondaryText)

	StyleCurrency = lipgloss.NewStyle().
			Foreground(ColorPrimaryText)

	StylePercentage = lipgloss.NewStyle().
			Foreground(ColorSecondaryText)
)

func StyleValue(positive bool) lipgloss.Style {
	if positive {
		return StyleProfit
	}
	return StyleLoss
}

func FormatProfit(value float64, currency string) string {
	if value >= 0 {
		return StyleProfit.Render("+" + currency + formatNumber(value))
	}
	return StyleLoss.Render("-" + currency + formatNumber(-value))
}

func formatNumber(val float64) string {
	return intToStr(int(val + 0.5))
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
