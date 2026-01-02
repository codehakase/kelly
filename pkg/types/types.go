package types

type CalculationMethod string

const (
	MethodArbitrage    CalculationMethod = "arbitrage"
	MethodKelly        CalculationMethod = "kelly"
	MethodProportional CalculationMethod = "proportional"
)

type OddsFormat string

const (
	FormatDecimal    OddsFormat = "decimal"
	FormatPercentage OddsFormat = "percentage"
	FormatFractional OddsFormat = "fractional"
	FormatAmerican   OddsFormat = "american"
)

type OutputFormat string

const (
	OutputTable OutputFormat = "table"
	OutputJSON  OutputFormat = "json"
	OutputCSV   OutputFormat = "csv"
)

type Option struct {
	Name               string  `json:"name"`
	Odds               float64 `json:"odds"`
	ImpliedProbability float64 `json:"implied_probability"`
	Probability        float64 `json:"probability,omitempty"`
	Stake              float64 `json:"stake"`
	ReturnIfWins       float64 `json:"return_if_wins"`
	ProfitIfWins       float64 `json:"profit_if_wins"`
	ROI                float64 `json:"roi"`
}

type Summary struct {
	GuaranteedProfit bool    `json:"guaranteed_profit"`
	MinProfit        float64 `json:"min_profit"`
	MaxProfit        float64 `json:"max_profit"`
	ExpectedValue    float64 `json:"expected_value"`
	MinROI           float64 `json:"min_roi"`
	MaxROI           float64 `json:"max_roi"`
	MarketEfficiency float64 `json:"market_efficiency"`
}

type CalculationResult struct {
	Method     CalculationMethod `json:"method"`
	TotalStake float64           `json:"total_stake"`
	Currency   string            `json:"currency"`
	OptionA    Option            `json:"option_a"`
	OptionB    Option            `json:"option_b"`
	Summary    Summary           `json:"summary"`
}

type CalculationInput struct {
	Method     CalculationMethod
	OddsA      float64
	OddsB      float64
	TotalStake float64
	ProbA      float64
	ProbB      float64
	NameA      string
	NameB      string
	Currency   string
}
