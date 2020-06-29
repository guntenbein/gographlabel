package gographlabel

type CurrentVertexLabelingRule struct {
	Name        string
	CurrentType string
	ResultLabel string
}

func (r CurrentVertexLabelingRule) ApplyRule(v *Vertex) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	v.Add(Label(r.ResultLabel))
	return true, nil
}

type ParentVertexLabelingRule struct {
	Name        string
	CurrentType string
	ParentType  string
	ResultLabel string
}

func (r ParentVertexLabelingRule) ApplyRule(v *Vertex) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	applied := false
	if err := v.ApplyParents(func(visitedVertex *Vertex) {
		if visitedVertex == nil {
			return
		}
		if r.ParentType == "" || visitedVertex.Type == r.ParentType {
			visitedVertex.Add(Label(r.ResultLabel))
			applied = true
		}
	}); err != nil {
		return false, err
	}
	return applied, nil
}

type ChildrenVertexLabelingRule struct {
	Name        string
	CurrentType string
	ChildType   string
	ResultLabel string
}

func (r ChildrenVertexLabelingRule) ApplyRule(v *Vertex) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	applied := false
	if err := v.ApplyChildren(func(visitedVertex *Vertex) {
		if visitedVertex == nil {
			return
		}
		if r.ChildType == "" || visitedVertex.Type == r.ChildType {
			visitedVertex.Add(Label(r.ResultLabel))
			applied = true
		}
	}); err != nil {
		return false, err
	}
	return applied, nil
}
