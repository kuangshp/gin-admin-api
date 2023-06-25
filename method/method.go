package method

type CommonMethod struct {
	ID   int32
	Name *string
}

func (m *CommonMethod) IsEmpty() bool {
	if m == nil {
		return true
	}
	return m.ID == 0
}

func (m *CommonMethod) GetName() string {
	if m == nil || m.Name == nil {
		return ""
	}
	return *m.Name
}
