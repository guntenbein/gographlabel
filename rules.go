package gographlabel

import "fmt"

type CurrentVertexLabelingRule struct {
	Name        string
	CurrentType string
	ResultLabel string
}

func (r CurrentVertexLabelingRule) ApplyRule(v *Vertex, cID string) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	if !v.Reserve(r.ResultLabel, cID) {
		return false, fmt.Errorf("vertex '%s' already labeled by '%s' for correlationID '%s'",
			v.ID, r.ResultLabel, v.MustGet(r.ResultLabel))
	}
	return true, nil
}

type ParentVertexLabelingRule struct {
	Name        string
	CurrentType string
	ParentType  string
	ResultLabel string
}

func (r ParentVertexLabelingRule) ApplyRule(v *Vertex, cID string) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	applied := false
	if err := v.ApplyParents(func(visitedVertex *Vertex) error {
		if visitedVertex == nil {
			return nil
		}
		if r.ParentType == "" || visitedVertex.Type == r.ParentType {
			if !visitedVertex.Reserve(r.ResultLabel, cID) {
				return fmt.Errorf("vertex '%s' already labeled by '%s' for correlationID '%s'",
					visitedVertex.ID, r.ResultLabel, v.MustGet(r.ResultLabel))
			}
			applied = true
		}
		return nil
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

func (r ChildrenVertexLabelingRule) ApplyRule(v *Vertex, cID string) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	applied := false
	if err := v.ApplyChildren(func(visitedVertex *Vertex) error {
		if visitedVertex == nil {
			return nil
		}
		if r.ChildType == "" || visitedVertex.Type == r.ChildType {
			if !visitedVertex.Reserve(r.ResultLabel, cID) {
				return fmt.Errorf("vertex '%s' already labeled by '%s' for correlationID '%s'",
					visitedVertex.ID, r.ResultLabel, v.MustGet(r.ResultLabel))
			}
			applied = true
		}
		return nil
	}); err != nil {
		return false, err
	}
	return applied, nil
}

type BrotherVertexLabelingRule struct {
	Name        string
	CurrentType string
	ParentType  string
	BrotherType string
	ResultLabel string
}

func (r BrotherVertexLabelingRule) ApplyRule(v *Vertex, cID string) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	// Find right parent
	var parent *Vertex
	if err := v.ApplyParents(func(visitedVertex *Vertex) error {
		if visitedVertex == nil || visitedVertex == v {
			return nil
		}
		if r.ParentType == "" || visitedVertex.Type == r.ParentType {
			parent = visitedVertex
			return nil
		}
		return nil
	}); err != nil {
		return false, err
	}
	if parent == nil {
		return false, nil
	}
	markChildrenRule := ChildrenVertexLabelingRule{"internal for BrotherVertexLabelingRule", r.ParentType, r.BrotherType, r.ResultLabel}
	return markChildrenRule.ApplyRule(parent, cID)
}
