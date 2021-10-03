package ken

type ObjectProvider interface {
	Get(key string) interface{}
}

type ObjectInjector interface {
	Set(key string, obj interface{})
}

type ObjectMap interface {
	ObjectProvider
	ObjectInjector

	Purge()
}

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
