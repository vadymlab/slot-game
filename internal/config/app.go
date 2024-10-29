package config

import "github.com/urfave/cli/v2"

// Constants for flag names used in SlotConfig
const (
	multiplierThree       = "multiplier-three"        // Flag for multiplier when three symbols match
	multiplierTwo         = "multiplier-two"          // Flag for multiplier when two symbols match
	twoMatchProbability   = "two-match-probability"   // Flag for probability of winning with two matches
	threeMatchProbability = "three-match-probability" // Flag for probability of winning with three matches
	rateLIMIT             = "rate-limit"              // Flag for rate limit (requests per second)
)

// SlotConfig defines configuration parameters for the slot game,
// including multipliers and probabilities for different winning scenarios.
type SlotConfig struct {
	MultiplierThree       float64 // Multiplier applied when three symbols match
	MultiplierTwo         float64 // Multiplier applied when two symbols match
	TwoMatchProbability   float64 // Probability for winning with two matching symbols
	ThreeMatchProbability float64 // Probability for winning with three matching symbols
	RateLimit             string  // Rate limit for requests per second
}

// GetSlotConfig returns a SlotConfig instance populated from CLI context flags.
// This function extracts values for each configuration parameter from the provided
// cli.Context, allowing dynamic configuration through command-line flags.
//
// Parameters:
//   - c: The CLI context from which to retrieve flag values.
//
// Returns:
//
//	A pointer to a SlotConfig struct with values obtained from the CLI flags.
func GetSlotConfig(c *cli.Context) *SlotConfig {
	return &SlotConfig{
		MultiplierThree:       c.Float64(multiplierThree),
		MultiplierTwo:         c.Float64(multiplierTwo),
		TwoMatchProbability:   c.Float64(twoMatchProbability),
		ThreeMatchProbability: c.Float64(threeMatchProbability),
		RateLimit:             c.String(rateLIMIT),
	}
}

// SlotFlags defines the command-line flags for configuring the slot game,
// including multipliers and probabilities for different win conditions.
// Each flag is linked to an environment variable, allowing for
// configuration via the environment as well as the CLI.
var SlotFlags = []cli.Flag{
	&cli.Float64Flag{
		Name:    multiplierThree,
		Value:   10,
		Usage:   "Multiplier for three matching symbols",
		EnvVars: []string{"MULTIPLIER_THREE"}, // Environment variable for multiplier on three matches
	},
	&cli.Float64Flag{
		Name:    multiplierTwo,
		Value:   2,
		Usage:   "Multiplier for two matching symbols",
		EnvVars: []string{"MULTIPLIER_TWO"}, // Environment variable for multiplier on two matches
	},
	&cli.Float64Flag{
		Name:    twoMatchProbability,
		Value:   0.30,
		Usage:   "Probability for winning with two matching symbols",
		EnvVars: []string{"TWO_MATCH_PROBABILITY"}, // Environment variable for probability on two matches
	},
	&cli.Float64Flag{
		Name:    threeMatchProbability,
		Value:   0.05,
		Usage:   "Probability for winning with three matching symbols",
		EnvVars: []string{"THREE_MATCH_PROBABILITY"}, // Environment variable for probability on three matches
	},
	&cli.StringFlag{
		Name:    rateLIMIT,
		Value:   "1-S",
		Usage:   "Rate limit for requests per second( 5 reqs/second: \"5-S\", 10 reqs/minute: \"10-M\", 100 reqs/hour: \"100-H\")",
		EnvVars: []string{"RATE_LIMIT"}, // Environment variable for rate limit
	},
}
