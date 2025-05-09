package ability

import "errors"

var (
	ErrAbilityHandlerNotDefined = errors.New("ability handler is not defined")
	ErrAbilityItemNotFound      = errors.New("ability item not found")
)
