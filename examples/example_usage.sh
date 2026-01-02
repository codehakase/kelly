#!/bin/bash
# Example usage scenarios for kelly

echo "=== Kelly Betting Calculator Examples ==="
echo

# Build the binary first
echo "Building kelly..."
go build -o kelly ..
echo

echo "=== Example 1: Grammy Awards Arbitrage ==="
echo "Command: kelly -a 39% -b 26% -t 10000 --name-a 'Davido' --name-b 'Tyla'"
./kelly -a 39% -b 26% -t 10000 --name-a "Davido - With You" --name-b "Tyla - PUSH 2 START"
echo

echo "=== Example 2: Basic Decimal Odds ==="
echo "Command: kelly -a 2.5 -b 3.0 -t 1000"
./kelly -a 2.5 -b 3.0 -t 1000
echo

echo "=== Example 3: Fractional Odds ==="
echo "Command: kelly -a 3/2 -b 5/2 -t 500"
./kelly -a 3/2 -b 5/2 -t 500
echo

echo "=== Example 4: American Odds ==="
echo "Command: kelly -a +150 -b +200 -t 1000"
./kelly -a +150 -b +200 -t 1000
echo

echo "=== Example 5: Kelly Criterion ==="
echo "Command: kelly -a 2.1 -b 3.5 -t 1000 --method kelly --prob-a 0.55 --prob-b 0.40"
./kelly -a 2.1 -b 3.5 -t 1000 --method kelly --prob-a 0.55 --prob-b 0.40
echo

echo "=== Example 6: Proportional Allocation ==="
echo "Command: kelly -a 2.5 -b 4.0 -t 1000 --method proportional"
./kelly -a 2.5 -b 4.0 -t 1000 --method proportional
echo

echo "=== Example 7: JSON Output ==="
echo "Command: kelly -a 2.56 -b 3.85 -t 10000 -f json"
./kelly -a 2.56 -b 3.85 -t 10000 -f json
echo

echo "=== Example 8: CSV Output ==="
echo "Command: kelly -a 2.56 -b 3.85 -t 10000 -f csv"
./kelly -a 2.56 -b 3.85 -t 10000 -f csv
echo

echo "=== Example 9: Compare All Methods ==="
echo "Command: kelly -a 2.5 -b 3.0 -t 1000 --compare --prob-a 0.5 --prob-b 0.4"
./kelly -a 2.5 -b 3.0 -t 1000 --compare --prob-a 0.5 --prob-b 0.4
echo

echo "=== Example 10: Verbose Output ==="
echo "Command: kelly -a 2.56 -b 3.85 -t 10000 -v"
./kelly -a 2.56 -b 3.85 -t 10000 -v
echo

# Cleanup
rm -f kelly
echo "Examples complete!"
