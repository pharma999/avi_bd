package game

import (
	"math"
	"math/rand"
	"time"

	"avi_bd/config"
)

// CrashGenerator generates crash points for game rounds
type CrashGenerator struct {
	cfg *config.Config
}

// NewCrashGenerator creates a new crash generator
func NewCrashGenerator(cfg *config.Config) *CrashGenerator {
	return &CrashGenerator{cfg: cfg}
}

// GenerateCrash generates a crash multiplier using a weighted distribution
// This ensures more realistic game distribution with lower crashes being more common
func (cg *CrashGenerator) GenerateCrash() float64 {
	// Seed random with current time for better randomness
	rand.Seed(time.Now().UnixNano())

	// Weighted distribution: exponential-ish distribution
	// More low crashes, fewer very high crashes
	r := rand.Float64()

	// Exponential distribution scaled to our range
	// This makes crashes under 2.0x more common (60%+)
	// and progressively fewer as multiplier increases
	crash := cg.cfg.CrashMin * math.Pow(r, -0.5)

	// Ensure within bounds
	if crash > cg.cfg.CrashMax {
		crash = cg.cfg.CrashMax
	}

	// Round to 2 decimals
	crash = math.Round(crash*100) / 100

	return crash
}

// GenerateCrashRange generates crash points within a specific range
func (cg *CrashGenerator) GenerateCrashRange(minCrash, maxCrash float64) float64 {
	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	crash := minCrash * math.Pow(r, -0.5)

	if crash > maxCrash {
		crash = maxCrash
	}

	return math.Round(crash*100) / 100
}
