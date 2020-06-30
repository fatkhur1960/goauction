package types

// ValidatorType query
type ValidatorType int

const (
	// ValidateJSON query type
	ValidateJSON ValidatorType = 1
	// ValidateQuery query type
	ValidateQuery ValidatorType = 2
	// ValidateURI query type
	ValidateURI ValidatorType = 3
)
