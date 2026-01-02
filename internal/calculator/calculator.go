package calculator

import (
	"errors"
	"math"

	"github.com/codehakase/kelly/pkg/types"
)

type Calculator interface {
	Calculate(input *types.CalculationInput) (*types.CalculationResult, error)
}

func NewCalculator(method types.CalculationMethod) Calculator {
	switch method {
	case types.MethodKelly:
		return &KellyCalculator{}
	case types.MethodProportional:
		return &ProportionalCalculator{}
	default:
		return &ArbitrageCalculator{}
	}
}

func round(val float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(val*multiplier) / multiplier
}

func impliedProbability(odds float64) float64 {
	if odds <= 0 {
		return 0
	}
	return 1.0 / odds
}

func marketEfficiency(oddsA, oddsB float64) float64 {
	return impliedProbability(oddsA) + impliedProbability(oddsB)
}

// ArbitrageCalculator implements guaranteed profit allocation.
type ArbitrageCalculator struct{}

func (c *ArbitrageCalculator) Calculate(input *types.CalculationInput) (*types.CalculationResult, error) {
	denominator := input.OddsA + input.OddsB - 2.0
	stakeA := round(input.TotalStake*(input.OddsB-1.0)/denominator, 2)
	stakeB := round(input.TotalStake*(input.OddsA-1.0)/denominator, 2)

	returnA := stakeA * input.OddsA
	returnB := stakeB * input.OddsB
	profitA := returnA - input.TotalStake
	profitB := returnB - input.TotalStake

	marketEff := marketEfficiency(input.OddsA, input.OddsB)

	return &types.CalculationResult{
		Method:     types.MethodArbitrage,
		TotalStake: input.TotalStake,
		Currency:   input.Currency,
		OptionA: types.Option{
			Name:               input.NameA,
			Odds:               input.OddsA,
			ImpliedProbability: impliedProbability(input.OddsA),
			Stake:              stakeA,
			ReturnIfWins:       round(returnA, 2),
			ProfitIfWins:       round(profitA, 2),
			ROI:                round(profitA/input.TotalStake, 4),
		},
		OptionB: types.Option{
			Name:               input.NameB,
			Odds:               input.OddsB,
			ImpliedProbability: impliedProbability(input.OddsB),
			Stake:              stakeB,
			ReturnIfWins:       round(returnB, 2),
			ProfitIfWins:       round(profitB, 2),
			ROI:                round(profitB/input.TotalStake, 4),
		},
		Summary: types.Summary{
			GuaranteedProfit: marketEff < 1.0,
			MinProfit:        round(math.Min(profitA, profitB), 2),
			MaxProfit:        round(math.Max(profitA, profitB), 2),
			ExpectedValue:    round((profitA+profitB)/2.0, 2),
			MinROI:           round(math.Min(profitA, profitB)/input.TotalStake, 4),
			MaxROI:           round(math.Max(profitA, profitB)/input.TotalStake, 4),
			MarketEfficiency: round(marketEff, 4),
		},
	}, nil
}

// KellyCalculator implements Kelly Criterion allocation.
type KellyCalculator struct{}

func (c *KellyCalculator) Calculate(input *types.CalculationInput) (*types.CalculationResult, error) {
	if input.ProbA == 0 || input.ProbB == 0 {
		return nil, errors.New("kelly method requires probability estimates for both options")
	}

	kellyA := math.Max(0, (input.ProbA*input.OddsA-1.0)/(input.OddsA-1.0))
	kellyB := math.Max(0, (input.ProbB*input.OddsB-1.0)/(input.OddsB-1.0))

	rawStakeA := input.TotalStake * kellyA
	rawStakeB := input.TotalStake * kellyB
	totalRaw := rawStakeA + rawStakeB

	var stakeA, stakeB float64
	if totalRaw > input.TotalStake {
		scale := input.TotalStake / totalRaw
		stakeA = round(rawStakeA*scale, 2)
		stakeB = round(rawStakeB*scale, 2)
	} else {
		stakeA = round(rawStakeA, 2)
		stakeB = round(rawStakeB, 2)
	}

	returnA := stakeA * input.OddsA
	returnB := stakeB * input.OddsB
	profitA := returnA - input.TotalStake
	profitB := returnB - input.TotalStake

	expectedValue := (input.ProbA * profitA) + (input.ProbB * profitB)
	if probSum := input.ProbA + input.ProbB; probSum < 1.0 {
		expectedValue += (1.0 - probSum) * (-input.TotalStake)
	}

	marketEff := marketEfficiency(input.OddsA, input.OddsB)

	return &types.CalculationResult{
		Method:     types.MethodKelly,
		TotalStake: input.TotalStake,
		Currency:   input.Currency,
		OptionA: types.Option{
			Name:               input.NameA,
			Odds:               input.OddsA,
			ImpliedProbability: impliedProbability(input.OddsA),
			Probability:        input.ProbA,
			Stake:              stakeA,
			ReturnIfWins:       round(returnA, 2),
			ProfitIfWins:       round(profitA, 2),
			ROI:                round(profitA/input.TotalStake, 4),
		},
		OptionB: types.Option{
			Name:               input.NameB,
			Odds:               input.OddsB,
			ImpliedProbability: impliedProbability(input.OddsB),
			Probability:        input.ProbB,
			Stake:              stakeB,
			ReturnIfWins:       round(returnB, 2),
			ProfitIfWins:       round(profitB, 2),
			ROI:                round(profitB/input.TotalStake, 4),
		},
		Summary: types.Summary{
			GuaranteedProfit: marketEff < 1.0,
			MinProfit:        round(math.Min(profitA, profitB), 2),
			MaxProfit:        round(math.Max(profitA, profitB), 2),
			ExpectedValue:    round(expectedValue, 2),
			MinROI:           round(math.Min(profitA, profitB)/input.TotalStake, 4),
			MaxROI:           round(math.Max(profitA, profitB)/input.TotalStake, 4),
			MarketEfficiency: round(marketEff, 4),
		},
	}, nil
}

// ProportionalCalculator implements proportional allocation.
type ProportionalCalculator struct{}

func (c *ProportionalCalculator) Calculate(input *types.CalculationInput) (*types.CalculationResult, error) {
	weightA := 1.0 / input.OddsA
	weightB := 1.0 / input.OddsB
	totalWeight := weightA + weightB

	stakeA := round(input.TotalStake*(weightA/totalWeight), 2)
	stakeB := round(input.TotalStake*(weightB/totalWeight), 2)

	returnA := stakeA * input.OddsA
	returnB := stakeB * input.OddsB
	profitA := returnA - input.TotalStake
	profitB := returnB - input.TotalStake

	marketEff := marketEfficiency(input.OddsA, input.OddsB)

	return &types.CalculationResult{
		Method:     types.MethodProportional,
		TotalStake: input.TotalStake,
		Currency:   input.Currency,
		OptionA: types.Option{
			Name:               input.NameA,
			Odds:               input.OddsA,
			ImpliedProbability: impliedProbability(input.OddsA),
			Stake:              stakeA,
			ReturnIfWins:       round(returnA, 2),
			ProfitIfWins:       round(profitA, 2),
			ROI:                round(profitA/input.TotalStake, 4),
		},
		OptionB: types.Option{
			Name:               input.NameB,
			Odds:               input.OddsB,
			ImpliedProbability: impliedProbability(input.OddsB),
			Stake:              stakeB,
			ReturnIfWins:       round(returnB, 2),
			ProfitIfWins:       round(profitB, 2),
			ROI:                round(profitB/input.TotalStake, 4),
		},
		Summary: types.Summary{
			GuaranteedProfit: marketEff < 1.0,
			MinProfit:        round(math.Min(profitA, profitB), 2),
			MaxProfit:        round(math.Max(profitA, profitB), 2),
			ExpectedValue:    round((profitA+profitB)/2.0, 2),
			MinROI:           round(math.Min(profitA, profitB)/input.TotalStake, 4),
			MaxROI:           round(math.Max(profitA, profitB)/input.TotalStake, 4),
			MarketEfficiency: round(marketEff, 4),
		},
	}, nil
}
