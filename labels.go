package gographlabel

type LabelEnum map[string]string

func (le LabelEnum) Get(l string) (string, bool) {
	cID, ok := le[l]
	return cID, ok
}

func (le LabelEnum) MustGet(l string) string {
	return le[l]
}

func (le LabelEnum) Reserve(label, correlationID string) bool {
	if value, ok := le.Get(label); ok {
		if value != correlationID {
			return false
		}
	}
	le[label] = correlationID
	return true
}
