# Kelly

**Optimal Betting Stake Calculator**

A professional terminal-based calculator for optimal bet allocation using Kelly Criterion, arbitrage, and proportional strategies. Features a Bloomberg Terminal-inspired design.

<img width="1051" height="816" alt="Screenshot 2026-01-02 at 16 38 38" src="https://github.com/user-attachments/assets/528f0a03-6baf-4f30-bbfc-1f4a3c2363b1" />

## Features

- **Three calculation methods**: Arbitrage (guaranteed profit), Kelly Criterion (growth optimization), Proportional (inverse odds)
- **Multiple odds formats**: Decimal (2.5), Percentage (39%), Fractional (3/2), American (+250)
- **Dual interface**: Interactive TUI and command-line modes
- **Export formats**: Table, JSON, CSV
- **Real-time validation**: Input validation with helpful error messages

## Installation

### Via Go

```bash
go install github.com/codehakase/kelly@latest
```

### From Source

```bash
git clone https://github.com/codehakase/kelly
cd kelly
make install
```

## Quick Start

### Interactive TUI Mode

Launch the interactive interface by running without arguments:

```bash
kelly
```

### CLI Mode

```bash
# Basic arbitrage calculation
kelly -a 2.56 -b 3.85 -t 10000

# With percentage odds
kelly --odds-a 39% --odds-b 26% --total 10000

# Named options with custom currency
kelly -a 2.56 -b 3.85 -t 10000 \
  --name-a "Davido - With You" \
  --name-b "Tyla - PUSH 2 START" \
  --currency "₦"

# Kelly Criterion with probability estimates
kelly -a 2.1 -b 3.5 -t 1000 \
  --method kelly \
  --prob-a 0.55 --prob-b 0.40

# Export to JSON
kelly -a 2.56 -b 3.85 -t 10000 -f json

# Compare all methods
kelly -a 2.56 -b 3.85 -t 10000 --compare
```

## Keyboard Shortcuts (TUI Mode)

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Navigate between fields |
| `Enter` | Calculate allocation |
| `m` | Cycle calculation method |
| `c` | Toggle compare mode |
| `?` | Show help overlay |
| `Ctrl+C` / `q` | Quit |

## CLI Reference

```
FLAGS:
  -a, --odds-a      Odds for Option A (required)
  -b, --odds-b      Odds for Option B (required)
  -t, --total       Total amount to allocate (required)
  -m, --method      Calculation method: arbitrage, kelly, proportional (default: arbitrage)
  -pa, --prob-a     Probability for Option A (required for Kelly)
  -pb, --prob-b     Probability for Option B (required for Kelly)
  -na, --name-a     Name/label for Option A (default: "Option A")
  -nb, --name-b     Name/label for Option B (default: "Option B")
  -c, --currency    Currency symbol (default: "₦")
  -f, --format      Output format: table, json, csv (default: table)
  -i               Force interactive TUI mode
  -v, --verbose     Verbose output with explanations
  --compare         Compare all calculation methods
  --no-color        Disable colored output
  --version         Show version information
```

## Odds Formats

| Format | Example | Description |
|--------|---------|-------------|
| Decimal | `2.5` | European format, total return per unit stake |
| Percentage | `39%` | Implied probability, converted to decimal |
| Fractional | `3/2` | UK format, profit per unit stake |
| American | `+150` | US format, positive for underdogs, negative for favorites |

## Calculation Methods

### Arbitrage (Default)

Calculates stakes to guarantee profit regardless of outcome. Works when market efficiency < 100%.

**Formula:**
```
Stake_A = Total × (Odds_B - 1) / (Odds_A + Odds_B - 2)
Stake_B = Total × (Odds_A - 1) / (Odds_A + Odds_B - 2)
```

### Kelly Criterion

Optimizes stake size based on your probability estimates to maximize long-term growth. Requires probability inputs.

**Formula:**
```
Kelly% = (p × odds - 1) / (odds - 1)
```

### Proportional

Simple allocation inversely proportional to odds. Lower odds receive higher stakes.

**Formula:**
```
Weight_A = 1 / Odds_A
Stake_A = Total × (Weight_A / Total_Weight)
```

## Example Output

```
╭────────────────────────────────────────────────────────────────────╮
│ KELLY • Stake Calculator               Method: ARBITRAGE          │
╰────────────────────────────────────────────────────────────────────╯

ALLOCATION BREAKDOWN

OPTION A • Davido - With You       OPTION B • Tyla - PUSH 2 START
Odds          2.56 (39.06%)        Odds          3.85 (25.97%)
Stake         ₦6,453 (64.53%)      Stake         ₦3,547 (35.47%)
Return        ₦16,520              Return        ₦13,656
Profit        +₦6,520              Profit        +₦3,656
ROI           +65.20%              ROI           +36.56%

SUMMARY
Total Invested        ₦10,000
Guaranteed Profit     YES
Profit Range          ₦3,656 - ₦6,520
ROI Range             36.56% - 65.20%
Market Efficiency     93.03% (Arbitrage opportunity)
```

## Development

```bash
# Run tests
make test

# Build binary
make build

# Install locally
make install

# Run linters
make lint

# Clean build artifacts
make clean
```

## License

MIT License - see [LICENSE](LICENSE) for details.
