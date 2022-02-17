package ken

// ResponsePolicy describes rules for context
// followups and responses.
type ResponsePolicy struct {
	// When set to true, the command response will
	// only be received by the sender of the command.
	//
	// This sets the `Ephemeral` flag of the `Context`
	// to true before any middleware is invoked. So,
	// you are able to modify the ephemeral flag either
	// in your middleware or directly in your command
	// logic, if you desire.
	Ephemeral bool
}

// ResponsePolicyCommand defines a command which
// provides a ResponsePolicy.
type ResponsePolicyCommand interface {
	ResponsePolicy() ResponsePolicy
}

// EphemeralCommand can be added to your command
// to make all command responses ephemeral.
// This means, that all responses to the command
// from the bot will only be received by the sender
// of the command.
type EphemeralCommand struct{}

var _ ResponsePolicyCommand = (*EphemeralCommand)(nil)

func (EphemeralCommand) ResponsePolicy() ResponsePolicy {
	return ResponsePolicy{
		Ephemeral: true,
	}
}
