package compile

type CompilationError struct {
	message string
}

func (ce *CompilationError) Error() string {
	return ce.message
}
