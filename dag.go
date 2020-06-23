package gographlabel

type Vertex struct {
	Type     string
	Parent   *Vertex
	Children []Vertex
	Labels   []string
}

type Rule struct {
	Name           string
	Up             bool
	ConditionLabel string
	ResultLabel    string
}

func (*Vertex) ApplyRule(rule Rule) bool {
	return true
}
