package game

import (
	"math"
)

// PayoutCalculator calculates payouts
type PayoutCalculator struct{}

// NewPayoutCalculator creates a new payout calculator
func NewPayoutCalculator() *PayoutCalculator {
	return &PayoutCalculator{}
}

// Calculate calculates payout based on bet amount and multiplier
func (pc *PayoutCalculator) Calculate(betAmount float64, multiplier float64) float64 {
	payout := betAmount * multiplier
	// Round to 2 decimals
	return math.Round(payout*100) / 100
}

// CalculateProfit calculates profit (payout - original bet)
func (pc *PayoutCalculator) CalculateProfit(betAmount float64, multiplier float64) float64 {
	payout := pc.Calculate(betAmount, multiplier)
	profit := payout - betAmount
	return math.Round(profit*100) / 100
}
