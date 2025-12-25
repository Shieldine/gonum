package transform

import "fmt"

// DegenerateInputError represents an error due to input data with variance
// below a threshold, which would cause numerical instability.
type DegenerateInputError float64

func (e DegenerateInputError) Error() string {
	return fmt.Sprintf("variance too low: %v", float64(e))
}
