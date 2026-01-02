package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/codehakase/kelly/internal/calculator"
	"github.com/codehakase/kelly/internal/formatter"
	"github.com/codehakase/kelly/internal/parser"
	"github.com/codehakase/kelly/internal/ui"
	"github.com/codehakase/kelly/internal/validator"
	"github.com/codehakase/kelly/pkg/types"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	var (
		oddsA       = flag.String("a", "", "Odds for Option A (required for CLI mode)")
		oddsB       = flag.String("b", "", "Odds for Option B (required for CLI mode)")
		total       = flag.Float64("t", 0, "Total amount to allocate (required for CLI mode)")
		method      = flag.String("m", "arbitrage", "Calculation method (arbitrage, kelly, proportional)")
		probA       = flag.Float64("pa", 0, "Probability for Option A (required for Kelly method)")
		probB       = flag.Float64("pb", 0, "Probability for Option B (required for Kelly method)")
		nameA       = flag.String("na", "Option A", "Name/label for Option A")
		nameB       = flag.String("nb", "Option B", "Name/label for Option B")
		currency    = flag.String("c", "₦", "Currency symbol")
		format      = flag.String("f", "table", "Output format (table, json, csv)")
		interactive = flag.Bool("i", false, "Force interactive TUI mode")
		verbose     = flag.Bool("v", false, "Verbose output with explanations")
		noColor     = flag.Bool("no-color", false, "Disable colored output")
		compare     = flag.Bool("compare", false, "Compare all calculation methods")
		version     = flag.Bool("version", false, "Show version information")
	)

	flag.StringVar(oddsA, "odds-a", "", "Odds for Option A")
	flag.StringVar(oddsB, "odds-b", "", "Odds for Option B")
	flag.Float64Var(total, "total", 0, "Total amount to allocate")
	flag.StringVar(method, "method", "arbitrage", "Calculation method")
	flag.Float64Var(probA, "prob-a", 0, "Probability for Option A")
	flag.Float64Var(probB, "prob-b", 0, "Probability for Option B")
	flag.StringVar(nameA, "name-a", "Option A", "Name for Option A")
	flag.StringVar(nameB, "name-b", "Option B", "Name for Option B")
	flag.BoolVar(verbose, "verbose", false, "Verbose output")

	flag.Usage = printUsage
	flag.Parse()

	if *version {
		fmt.Printf("Kelly Calculator %s (built %s)\n", Version, BuildTime)
		os.Exit(0)
	}

	if len(os.Args) == 1 || *interactive {
		runInteractive()
	} else if *oddsA != "" && *oddsB != "" && *total > 0 {
		runCLI(*oddsA, *oddsB, *total, *method, *probA, *probB,
			*nameA, *nameB, *currency, *format, *verbose, *noColor, *compare)
	} else {
		if *oddsA != "" || *oddsB != "" || *total > 0 {
			fmt.Fprintln(os.Stderr, "Error: CLI mode requires --odds-a, --odds-b, and --total")
			fmt.Fprintln(os.Stderr, "Run with -h for usage information")
			os.Exit(1)
		}
		runInteractive()
	}
}

func runInteractive() {
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCLI(oddsAStr, oddsBStr string, total float64, methodStr string,
	probA, probB float64, nameA, nameB, currency, format string,
	verbose, noColor, compare bool) {

	decimalOddsA, err := parser.ParseOdds(oddsAStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error parsing odds A: %v\n", err)
		os.Exit(1)
	}

	decimalOddsB, err := parser.ParseOdds(oddsBStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error parsing odds B: %v\n", err)
		os.Exit(1)
	}

	calcMethod := types.CalculationMethod(methodStr)
	switch calcMethod {
	case types.MethodArbitrage, types.MethodKelly, types.MethodProportional:
	default:
		fmt.Fprintf(os.Stderr, "✗ Error: Invalid method '%s'. Must be: arbitrage, kelly, or proportional\n", methodStr)
		os.Exit(1)
	}

	input := &types.CalculationInput{
		Method: calcMethod, OddsA: decimalOddsA, OddsB: decimalOddsB, TotalStake: total,
		ProbA: probA, ProbB: probB, NameA: nameA, NameB: nameB, Currency: currency,
	}

	if err := validator.ValidateCalculationInput(input); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Validation error: %v\n", err)
		os.Exit(1)
	}

	if compare {
		runComparison(input, format, verbose)
		return
	}

	calc := calculator.NewCalculator(input.Method)
	result, err := calc.Calculate(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Calculation error: %v\n", err)
		os.Exit(1)
	}

	output, err := formatOutput(result, format, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Formatting error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(output)
}

func runComparison(input *types.CalculationInput, format string, verbose bool) {
	methods := []types.CalculationMethod{
		types.MethodArbitrage, types.MethodKelly, types.MethodProportional,
	}

	fmt.Println("╭─────────────────────────────────────────────────────────────────────╮")
	fmt.Println("│ KELLY • Method Comparison                                           │")
	fmt.Println("╰─────────────────────────────────────────────────────────────────────╯")
	fmt.Println()

	for _, method := range methods {
		input.Method = method

		if method == types.MethodKelly && (input.ProbA == 0 || input.ProbB == 0) {
			fmt.Printf("─── %s (skipped: requires probabilities) ───\n\n", methodName(method))
			continue
		}

		calc := calculator.NewCalculator(method)
		result, err := calc.Calculate(input)
		if err != nil {
			fmt.Printf("─── %s (error: %v) ───\n\n", methodName(method), err)
			continue
		}

		fmt.Printf("─── %s ───\n", methodName(method))
		output, _ := formatOutput(result, format, verbose)
		fmt.Println(output)
		fmt.Println()
	}
}

func formatOutput(result *types.CalculationResult, format string, verbose bool) (string, error) {
	switch types.OutputFormat(format) {
	case types.OutputJSON:
		return formatter.FormatJSON(result)
	case types.OutputCSV:
		return formatter.FormatCSV(result)
	default:
		return formatter.FormatTable(result, verbose), nil
	}
}

func methodName(method types.CalculationMethod) string {
	switch method {
	case types.MethodArbitrage:
		return "ARBITRAGE (Guaranteed Profit)"
	case types.MethodKelly:
		return "KELLY CRITERION (Growth Optimization)"
	case types.MethodProportional:
		return "PROPORTIONAL (Inverse Odds)"
	default:
		return string(method)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Kelly - Optimal Betting Stake Calculator

USAGE:
  kelly                          Launch interactive TUI (default)
  kelly [flags]                  Run calculation with CLI arguments

EXAMPLES:
  kelly
  kelly -a 2.56 -b 3.85 -t 10000
  kelly --odds-a 39%% --odds-b 26%% --total 10000
  kelly -a 2.56 -b 3.85 -t 10000 --name-a "Davido" --name-b "Tyla" --currency "₦"
  kelly -a 2.1 -b 3.5 -t 1000 --method kelly --prob-a 0.55 --prob-b 0.40
  kelly -a 2.56 -b 3.85 -t 10000 -f json
  kelly -a 2.56 -b 3.85 -t 10000 --compare

FLAGS:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
ODDS FORMATS:
  Decimal:     2.5, 3.85
  Percentage:  39%%, 26%%
  Fractional:  3/2, 5/2
  American:    +250, -150

CALCULATION METHODS:
  arbitrage     Guarantees profit regardless of outcome (default)
  kelly         Maximizes growth based on probability estimates
  proportional  Simple allocation inversely proportional to odds

For more information, visit: https://github.com/codehakase/kelly
`)
}
