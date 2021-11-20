package ken

// UserCommand defines a callable user command.
type UserCommand interface {
	Command

	// TypeUser is used to differenciate between
	// UserCommand and MessageCommand which have
	// the same structure otherwise.
	//
	// This method must only be implemented and
	// will never be called by ken, so it can be
	// completely empty.
	TypeUser()
}
