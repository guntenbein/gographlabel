package gographlabel

type Label struct {
	Name string
}

type LabelEnum map[Label]struct{}

func (le LabelEnum) Contains(l Label) bool {
	_, ok := le[l]
	return ok
}

type Vertex struct {
	Type     string
	Parent   *Vertex
	Children []*Vertex
	Labels   LabelEnum
}

func (v *Vertex) Contains(label string) bool {
	v.Labels.Contains(Label{label})
	return false
}
