package bean

const Success = "Success"
const UnexpectedError = "UnexpectedError"

var CodeMessage = map[string]struct {
	Code    int
	Message string
}{
	Success: {1, "Success"},

	// -x for basic message
	UnexpectedError: {-1, "Unexpected error"},
}
