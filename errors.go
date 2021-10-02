package ken

import "errors"

var (
	ErrEmptyCommandName         = errors.New("command name can not be empty")
	ErrCommandAlreadyRegistered = errors.New("command with the same name has already been rgistered")
)
