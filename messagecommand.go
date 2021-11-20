package ken

// MessageCommand defines a callable message command.
type MessageCommand interface {
	Command

	// TypeMessage is used to differenciate between
	// UserCommand and MessageCommand which have
	// the same structure otherwise.
	//
	// This method must only be implemented and
	// will never be called by ken, so it can be
	// completely empty.
	TypeMessage()
}
