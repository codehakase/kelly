package ui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil
	}
	return m, m.updateActiveInput(msg)
}

func (m Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showHelp {
		switch msg.String() {
		case "?", "esc", "q", "enter":
			m.showHelp = false
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if m.getInputField(m.activeField).Value() == "" || !m.getInputField(m.activeField).Focused() {
			return m, tea.Quit
		}
		return m.updateInputAndRecalculate(msg)
	case "tab":
		return m, m.nextField()
	case "shift+tab":
		return m, m.prevField()
	case "enter":
		m.calculate()
		return m, nil
	case "m":
		if !m.isTypingLetter() {
			m.cycleMethod()
			return m, nil
		}
		return m.updateInputAndRecalculate(msg)
	case "c":
		if !m.isTypingLetter() {
			m.compareMode = !m.compareMode
			return m, nil
		}
		return m.updateInputAndRecalculate(msg)
	case "r":
		if !m.isTypingLetter() {
			m.reset()
			return m, nil
		}
		return m.updateInputAndRecalculate(msg)
	case "?":
		m.showHelp = true
		return m, nil
	default:
		return m.updateInputAndRecalculate(msg)
	}
}

func (m Model) updateInputAndRecalculate(msg tea.Msg) (Model, tea.Cmd) {
	cmd := m.updateActiveInput(msg)
	m.calculate()
	return m, cmd
}

func (m *Model) updateActiveInput(msg tea.Msg) tea.Cmd {
	return m.getInputField(m.activeField).Update(msg)
}

func (m Model) isTypingLetter() bool {
	field := m.getInputField(m.activeField)
	return field.Focused() && field.Value() != ""
}
