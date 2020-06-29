package gographlabel

type LabelEnum map[Label]struct{}

func (le LabelEnum) Contains(l Label) bool {
	_, ok := le[l]
	return ok
}

func (le LabelEnum) Add(l Label) bool {
	if le.Contains(l) {
		return false
	}
	le[l] = struct{}{}
	return true
}
