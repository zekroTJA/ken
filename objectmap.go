package ken

type ObjectProvider interface {
	Load(key interface{}) interface{}
}

type ObjectInjector interface {
	Store(key, obj interface{})
}

type ObjectMap interface {
	ObjectProvider
	ObjectInjector

	Purge()
}

type simpleObjectMap map[interface{}]interface{}

var _ ObjectMap = (simpleObjectMap)(nil)

func (m simpleObjectMap) Load(key interface{}) interface{} {
	return m[key]
}

func (m simpleObjectMap) Store(key, obj interface{}) {
	m[key] = obj
}

func (m simpleObjectMap) Purge() {
	for k := range m {
		delete(m, k)
	}
}
