package ken

// ObjectProvider specifies an instance providing
// objects by string key.
type ObjectProvider interface {
	// Get returns a stored object by its
	// key, if existent.
	Get(key string) interface{}
}

// ObjectInjector specifies an instance which
// allows storing an object by string key.
type ObjectInjector interface {
	// Set stores the given object by given
	// key.
	Set(key string, obj interface{})
}

// ObjectMap combines ObjectProvider and
// ObjectInjector.
type ObjectMap interface {
	ObjectProvider
	ObjectInjector

	// Purge cleans all stored objects and
	// keys from the provider.
	Purge()
}

// simpleObjectMap implements ObjectMap for
// a map[string]interface{}.
type simpleObjectMap map[string]interface{}

var _ ObjectMap = (simpleObjectMap)(nil)

func (m simpleObjectMap) Get(key string) interface{} {
	return m[key]
}

func (m simpleObjectMap) Set(key string, obj interface{}) {
	m[key] = obj
}

func (m simpleObjectMap) Purge() {
	for k := range m {
		delete(m, k)
	}
}
