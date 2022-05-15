package util

type Container interface {
	Contain(key string) bool
}

func IsValidSecret(c Container, all [][]string) (bool, []string) {
	for _, must := range all {
		if containAllKey(c, must) {
			return true, must
		}
	}
	return false, nil
}
func containAllKey(c Container, must []string) bool {
	for _, v := range must {
		if !c.Contain(v) {
			return false
		}
	}
	return true
}

type ssMap struct {
	m map[string]string
}

func (m *ssMap) Contain(key string) bool {
	_, exist := m.m[key]
	return exist
}

func WrapSSM(m map[string]string) Container {
	return &ssMap{
		m: m,
	}
}

/*func ContainsAllKey(data map[string]string, must []string) bool {
	return ContainAllKey(WrapSSM(data), must)
}
*/
type saMap struct {
	m map[string]interface{}
}

func (m *saMap) Contain(key string) bool {
	_, exist := m.m[key]
	return exist
}

func WrapSAM(m map[string]interface{}) Container {
	return &saMap{
		m: m,
	}
}

type sbMap struct {
	m map[string][]byte
}

func (m *sbMap) Contain(key string) bool {
	_, exist := m.m[key]
	return exist
}

func WrapSBM(m map[string][]byte) Container {
	return &sbMap{
		m: m,
	}
}
