package store

import (
	"encoding/json"
	"os"
)

// LocalCommandStore implements CommandStore for a
// local file as storage device.
type LocalCommandStore struct {
	loc string
}

var _ CommandStore = (*LocalCommandStore)(nil)

// NewLocalCommandStore creates a new instance of
// LocalCommandStore with the passed file location
// as stoage destination.
func NewLocalCommandStore(loc string) *LocalCommandStore {
	return &LocalCommandStore{loc}
}

// NewDefault returns a new LocalCommandStore with
// default file location (".commandCache.json").
func NewDefault() *LocalCommandStore {
	return NewLocalCommandStore(".commandCache.json")
}

func (lcs *LocalCommandStore) Store(cmds map[string]string) (err error) {
	f, err := os.Create(lcs.loc)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(cmds)
	return
}

func (lcs *LocalCommandStore) Load() (cmds map[string]string, err error) {
	cmds = map[string]string{}
	f, err := os.Open(lcs.loc)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&cmds)
	return
}
