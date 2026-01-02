package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/codehakase/kelly/internal/calculator"
	"github.com/codehakase/kelly/internal/parser"
	"github.com/codehakase/kelly/internal/ui/components"
	"github.com/codehakase/kelly/pkg/types"
)

const (
	fieldOddsA = iota
	fieldOddsB
	fieldTotal
	fieldNameA
	fieldNameB
	fieldProbA
	fieldProbB
	fieldCount
)

type Model struct {
	oddsAInput, oddsBInput, totalInput components.ValidatedInput
	nameAInput, nameBInput             components.ValidatedInput
	probAInput, probBInput             components.ValidatedInput

	activeField int
	method      types.CalculationMethod
	currency    string
	result      *types.CalculationResult
	err         error

	width, height int
	showHelp      bool
	compareMode   bool
	ready         bool
}

func NewModel() Model {
	m := Model{method: types.MethodArbitrage, currency: "₦"}

	m.oddsAInput = components.NewValidatedInput("Odds A", "2.56 or 39% or 3/2", validateOdds)
	m.oddsBInput = components.NewValidatedInput("Odds B", "3.85 or 26% or 5/2", validateOdds)
	m.totalInput = components.NewValidatedInput("Total", "10000", validateTotal)
	m.nameAInput = components.NewValidatedInput("Name A", "Option A", nil)
	m.nameBInput = components.NewValidatedInput("Name B", "Option B", nil)
	m.probAInput = components.NewValidatedInput("Prob A", "0.55 (for Kelly)", validateProbability)
	m.probBInput = components.NewValidatedInput("Prob B", "0.40 (for Kelly)", validateProbability)

	m.nameAInput.SetValue("Option A")
	m.nameBInput.SetValue("Option B")
	m.oddsAInput.Focus()

	return m
}

func (m Model) Init() tea.Cmd { return nil }

func validateOdds(input string) error {
	odds, err := parser.ParseOdds(input)
	if err != nil {
		return err
	}
	if odds < 1.01 {
		return fmt.Errorf("odds must be >= 1.01")
	}
	return nil
}

func validateTotal(input string) error {
	var total float64
	if _, err := fmt.Sscanf(input, "%f", &total); err != nil {
		return fmt.Errorf("invalid number")
	}
	if total <= 0 {
		return fmt.Errorf("must be positive")
	}
	return nil
}

func validateProbability(input string) error {
	var prob float64
	if _, err := fmt.Sscanf(input, "%f", &prob); err != nil {
		return fmt.Errorf("invalid number")
	}
	if prob <= 0 || prob >= 1 {
		return fmt.Errorf("must be between 0 and 1")
	}
	return nil
}

func (m *Model) getInputField(idx int) *components.ValidatedInput {
	switch idx {
	case fieldOddsA:
		return &m.oddsAInput
	case fieldOddsB:
		return &m.oddsBInput
	case fieldTotal:
		return &m.totalInput
	case fieldNameA:
		return &m.nameAInput
	case fieldNameB:
		return &m.nameBInput
	case fieldProbA:
		return &m.probAInput
	case fieldProbB:
		return &m.probBInput
	default:
		return &m.oddsAInput
	}
}

func (m *Model) focusField(idx int) tea.Cmd {
	m.oddsAInput.Blur()
	m.oddsBInput.Blur()
	m.totalInput.Blur()
	m.nameAInput.Blur()
	m.nameBInput.Blur()
	m.probAInput.Blur()
	m.probBInput.Blur()
	m.activeField = idx
	return m.getInputField(idx).Focus()
}

func (m *Model) nextField() tea.Cmd {
	next := m.activeField + 1
	if m.method != types.MethodKelly && (next == fieldProbA || next == fieldProbB) {
		next = fieldOddsA
	}
	if next >= fieldCount {
		next = fieldOddsA
	}
	return m.focusField(next)
}

func (m *Model) prevField() tea.Cmd {
	prev := m.activeField - 1
	if prev < 0 {
		if m.method == types.MethodKelly {
			prev = fieldProbB
		} else {
			prev = fieldNameB
		}
	}
	if m.method != types.MethodKelly && (prev == fieldProbA || prev == fieldProbB) {
		prev = fieldNameB
	}
	return m.focusField(prev)
}

func (m *Model) cycleMethod() {
	switch m.method {
	case types.MethodArbitrage:
		m.method = types.MethodKelly
	case types.MethodKelly:
		m.method = types.MethodProportional
	case types.MethodProportional:
		m.method = types.MethodArbitrage
	}
	m.calculate()
}

func (m *Model) calculate() {
	m.result = nil
	m.err = nil

	if !m.oddsAInput.IsValid() || !m.oddsBInput.IsValid() || !m.totalInput.IsValid() {
		return
	}

	oddsA, err := parser.ParseOdds(m.oddsAInput.Value())
	if err != nil {
		m.err = err
		return
	}

	oddsB, err := parser.ParseOdds(m.oddsBInput.Value())
	if err != nil {
		m.err = err
		return
	}

	var total float64
	if _, err = fmt.Sscanf(m.totalInput.Value(), "%f", &total); err != nil {
		m.err = fmt.Errorf("invalid total: %w", err)
		return
	}

	var probA, probB float64
	if m.method == types.MethodKelly {
		if m.probAInput.Value() != "" {
			fmt.Sscanf(m.probAInput.Value(), "%f", &probA)
		}
		if m.probBInput.Value() != "" {
			fmt.Sscanf(m.probBInput.Value(), "%f", &probB)
		}
		if probA == 0 || probB == 0 {
			m.err = fmt.Errorf("Kelly method requires probability estimates")
			return
		}
	}

	nameA := m.nameAInput.Value()
	if nameA == "" {
		nameA = "Option A"
	}
	nameB := m.nameBInput.Value()
	if nameB == "" {
		nameB = "Option B"
	}

	input := &types.CalculationInput{
		Method: m.method, OddsA: oddsA, OddsB: oddsB, TotalStake: total,
		ProbA: probA, ProbB: probB, NameA: nameA, NameB: nameB, Currency: m.currency,
	}

	calc := calculator.NewCalculator(m.method)
	result, err := calc.Calculate(input)
	if err != nil {
		m.err = err
		return
	}
	m.result = result
}

func (m *Model) reset() {
	m.oddsAInput.Reset()
	m.oddsBInput.Reset()
	m.totalInput.Reset()
	m.nameAInput.SetValue("Option A")
	m.nameBInput.SetValue("Option B")
	m.probAInput.Reset()
	m.probBInput.Reset()
	m.result = nil
	m.err = nil
	m.focusField(fieldOddsA)
}
