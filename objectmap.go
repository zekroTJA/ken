package ken

type ObjectProvider interface {
	Load(key string) interface{}
}

type ObjectInjector interface {
	Store(key string, obj interface{})
}

type ObjectMap interface {
	ObjectProvider
	ObjectInjector

	Purge()
}

type simpleObjectMap map[string]interface{}

var _ ObjectMap = (simpleObjectMap)(nil)

func (m simpleObjectMap) Load(key string) interface{} {
	return m[key]
}

func (m simpleObjectMap) Store(key string, obj interface{}) {
	m[key] = obj
}

func (m simpleObjectMap) Purge() {
	for k := range m {
		delete(m, k)
	}
}
