package store

// CommandStore allows to store and load registered
// commands so that they can be updated instead
// of deleted and re-created on restarts.
type CommandStore interface {
	// Store stores the passed commands map.
	Store(cmds map[string]string) error
	// Load retrieves a stored commands map or
	// an empty map, if no store was applied
	// before.
	Load() (cmds map[string]string, err error)
}
