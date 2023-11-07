package types

type ValidationInput struct {
	Key          string
	Value        interface{}
	Item         map[string]interface{}
	ExistingItem map[string]interface{}
}

type ValidationOutput struct {
	Input   *ValidationInput
	Result  bool
	Message string
	Error   error
}

type FieldValidation func(input *ValidationInput, ch chan ValidationOutput)
