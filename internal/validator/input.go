package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codehakase/kelly/pkg/types"
)

type ValidationError struct {
	Errors []error
}

func (e ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	var sb strings.Builder
	sb.WriteString("multiple validation errors:\n")
	for i, err := range e.Errors {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}
	return sb.String()
}

func ValidateOdds(odds float64) error {
	if odds < 1.01 {
		return fmt.Errorf("odds must be >= 1.01, got: %.2f", odds)
	}
	return nil
}

func ValidateProbability(prob float64) error {
	if prob <= 0 || prob >= 1 {
		return fmt.Errorf("probability must be between 0 and 1 (exclusive), got: %.4f", prob)
	}
	return nil
}

func ValidateTotalStake(total float64) error {
	if total <= 0 {
		return fmt.Errorf("total stake must be positive, got: %.2f", total)
	}
	return nil
}

func ValidateCalculationInput(input *types.CalculationInput) error {
	var errs []error

	if err := ValidateOdds(input.OddsA); err != nil {
		errs = append(errs, fmt.Errorf("Option A: %w", err))
	}
	if err := ValidateOdds(input.OddsB); err != nil {
		errs = append(errs, fmt.Errorf("Option B: %w", err))
	}
	if err := ValidateTotalStake(input.TotalStake); err != nil {
		errs = append(errs, err)
	}

	switch input.Method {
	case types.MethodKelly:
		if input.ProbA == 0 || input.ProbB == 0 {
			errs = append(errs, errors.New("Kelly method requires probability estimates for both options (use --prob-a and --prob-b)"))
		}
		if input.ProbA != 0 {
			if err := ValidateProbability(input.ProbA); err != nil {
				errs = append(errs, fmt.Errorf("Option A probability: %w", err))
			}
		}
		if input.ProbB != 0 {
			if err := ValidateProbability(input.ProbB); err != nil {
				errs = append(errs, fmt.Errorf("Option B probability: %w", err))
			}
		}
		if input.ProbA > 0 && input.ProbB > 0 {
			if sum := input.ProbA + input.ProbB; sum > 1.0 {
				errs = append(errs, fmt.Errorf("warning: probabilities sum to %.4f (> 1.0)", sum))
			}
		}
	case types.MethodArbitrage, types.MethodProportional:
		// No probability requirements
	default:
		errs = append(errs, fmt.Errorf("invalid calculation method: %s", input.Method))
	}

	if input.OddsA > 0 && input.OddsB > 0 {
		marketEff := (1.0 / input.OddsA) + (1.0 / input.OddsB)
		if input.Method == types.MethodArbitrage && marketEff >= 1.0 {
			errs = append(errs, fmt.Errorf("warning: combined implied probability (%.2f%%) >= 100%% - no guaranteed profit", marketEff*100))
		}
	}

	if len(errs) > 0 {
		return ValidationError{Errors: errs}
	}
	return nil
}

func ValidateCalculationInputStrict(input *types.CalculationInput) error {
	return ValidateCalculationInput(input)
}
