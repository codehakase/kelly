package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorMuted     = lipgloss.Color("#6b7280")
	colorAccent    = lipgloss.Color("#60a5fa")
	colorPrimary   = lipgloss.Color("#e4e4e7")
	colorSecondary = lipgloss.Color("#9ca3af")
	colorPanelBG   = lipgloss.Color("#1a1e3f")
	colorBorder    = lipgloss.Color("#2d3561")
	colorError     = lipgloss.Color("#ef4444")
)

var (
	inputLabelStyle    = lipgloss.NewStyle().Foreground(colorSecondary).Width(12)
	inputActiveStyle   = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	inputInactiveStyle = lipgloss.NewStyle().Foreground(colorPrimary)
	inputErrorStyle    = lipgloss.NewStyle().Foreground(colorError).Italic(true)
	inputCursorStyle   = lipgloss.NewStyle().Foreground(colorAccent)

	panelStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorBorder).Padding(1, 2)
	panelTitleStyle = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)

	helpKeyStyle  = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	helpDescStyle = lipgloss.NewStyle().Foreground(colorMuted)
	helpSepStyle  = lipgloss.NewStyle().Foreground(colorSecondary)
)

type ValidatedInput struct {
	Input     textinput.Model
	Label     string
	Validator func(string) error
	Error     error
	focused   bool
}

func NewValidatedInput(label, placeholder string, validator func(string) error) ValidatedInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 50
	ti.Width = 20
	ti.PromptStyle = inputActiveStyle
	ti.TextStyle = inputInactiveStyle
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorMuted)
	ti.Cursor.Style = inputCursorStyle

	return ValidatedInput{Input: ti, Label: label, Validator: validator}
}

func (vi *ValidatedInput) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	vi.Input, cmd = vi.Input.Update(msg)

	if vi.Validator != nil && vi.Input.Value() != "" {
		vi.Error = vi.Validator(vi.Input.Value())
	} else {
		vi.Error = nil
	}
	return cmd
}

func (vi ValidatedInput) View() string {
	labelStyle := inputLabelStyle
	if vi.focused {
		labelStyle = inputLabelStyle.Foreground(colorAccent).Bold(true)
	}

	result := labelStyle.Render(vi.Label) + " " + vi.Input.View()
	if vi.Error != nil {
		result += "\n             " + inputErrorStyle.Render(" ✗ "+vi.Error.Error())
	}
	return result
}

func (vi ValidatedInput) Value() string { return vi.Input.Value() }
func (vi ValidatedInput) Focused() bool { return vi.focused }
func (vi ValidatedInput) IsValid() bool { return vi.Input.Value() != "" && vi.Error == nil }
func (vi *ValidatedInput) SetValue(v string) {
	vi.Input.SetValue(v)
	if vi.Validator != nil && v != "" {
		vi.Error = vi.Validator(v)
	} else {
		vi.Error = nil
	}
}
func (vi *ValidatedInput) Focus() tea.Cmd { vi.focused = true; return vi.Input.Focus() }
func (vi *ValidatedInput) Blur()          { vi.focused = false; vi.Input.Blur() }
func (vi *ValidatedInput) Reset()         { vi.Input.SetValue(""); vi.Error = nil }

func Panel(title, content string, width int) string {
	style := panelStyle.Width(width - 4)
	if title != "" {
		content = panelTitleStyle.Render(title) + "\n" + content
	}
	return style.Render(content)
}

func PanelWithHeader(header, content string, width int) string {
	headerStyle := lipgloss.NewStyle().
		Foreground(colorPrimary).Background(colorPanelBG).Bold(true).Padding(0, 1).Width(width - 4)
	contentStyle := lipgloss.NewStyle().Padding(1, 0)
	return panelStyle.Width(width - 4).Render(headerStyle.Render(header) + "\n" + contentStyle.Render(content))
}

func HorizontalPanels(leftTitle, leftContent, rightTitle, rightContent string, totalWidth int) string {
	panelWidth := (totalWidth - 3) / 2
	return lipgloss.JoinHorizontal(lipgloss.Top,
		Panel(leftTitle, leftContent, panelWidth), " ",
		Panel(rightTitle, rightContent, panelWidth))
}

func SplitPanel(title, leftContent, rightContent string, width int) string {
	colWidth := (width - 8) / 2
	leftStyle := lipgloss.NewStyle().Width(colWidth)
	rightStyle := lipgloss.NewStyle().Width(colWidth)
	content := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftContent), "  ", rightStyle.Render(rightContent))
	return Panel(title, content, width)
}

func Help() string {
	keys := []struct{ key, desc string }{
		{"Tab", "Switch"}, {"Enter", "Calculate"}, {"m", "Method"},
		{"c", "Compare"}, {"?", "Help"}, {"q", "Quit"},
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, helpKeyStyle.Render("["+k.key+"]")+" "+helpDescStyle.Render(k.desc))
	}
	return helpSepStyle.Render(strings.Join(parts, "  "))
}

func HelpOverlay(width, height int) string {
	titleStyle := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Align(lipgloss.Center)
	sectionStyle := lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	keyStyle := lipgloss.NewStyle().Foreground(colorAccent).Width(15)
	descStyle := lipgloss.NewStyle().Foreground(colorPrimary)

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("KELLY CALCULATOR - HELP") + "\n\n")

	sb.WriteString(sectionStyle.Render("Navigation") + "\n")
	sb.WriteString(keyStyle.Render("Tab") + descStyle.Render("Move to next field") + "\n")
	sb.WriteString(keyStyle.Render("Shift+Tab") + descStyle.Render("Move to previous field") + "\n\n")

	sb.WriteString(sectionStyle.Render("Actions") + "\n")
	sb.WriteString(keyStyle.Render("Enter") + descStyle.Render("Calculate allocation") + "\n")
	sb.WriteString(keyStyle.Render("m") + descStyle.Render("Cycle calculation method") + "\n")
	sb.WriteString(keyStyle.Render("c") + descStyle.Render("Toggle comparison mode") + "\n")
	sb.WriteString(keyStyle.Render("r") + descStyle.Render("Reset all inputs") + "\n\n")

	sb.WriteString(sectionStyle.Render("General") + "\n")
	sb.WriteString(keyStyle.Render("?") + descStyle.Render("Toggle help") + "\n")
	sb.WriteString(keyStyle.Render("q / Ctrl+C") + descStyle.Render("Quit") + "\n\n")

	sb.WriteString(sectionStyle.Render("Methods") + "\n")
	sb.WriteString(descStyle.Render("• Arbitrage: Guaranteed profit\n"))
	sb.WriteString(descStyle.Render("• Kelly: Growth optimization\n"))
	sb.WriteString(descStyle.Render("• Proportional: Simple inverse allocation\n\n"))

	sb.WriteString(helpDescStyle.Render("Press ? or Esc to close"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).BorderForeground(colorBorder).
		Background(colorPanelBG).Padding(2, 4).Width(60)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, boxStyle.Render(sb.String()))
}
