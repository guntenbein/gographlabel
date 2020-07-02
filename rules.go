package gographlabel

import "fmt"

type CurrentVertexTypeCheckingRule struct {
	Name                string
	AllowedCurrentTypes []string
}

func (r CurrentVertexTypeCheckingRule) ApplyRule(v *Vertex, _ string) (bool, error) {
	if len(r.AllowedCurrentTypes) == 0 {
		return false, fmt.Errorf("rule '%s' should have at least one allowed type", r.Name)
	}
	for _, allowedType := range r.AllowedCurrentTypes {
		if allowedType == v.Type {
			return true, nil
		}
	}
	return false, fmt.Errorf("vertex '%s' does not have allowed type to apply the block, allowed types '%v', vertext type '%s'", v.ID, r.AllowedCurrentTypes, v.Type)
}

type CurrentVertexLabelingRule struct {
	Name        string
	CurrentType string
	ResultLabel string
	Exclusive   bool
}

func (r CurrentVertexLabelingRule) ApplyRule(v *Vertex, cID string) (bool, error) {
	if r.CurrentType != "" && r.CurrentType != v.Type {
		return false, nil
	}
	if !v.ReserveBlock(r.ResultLabel, cID, r.Exclusive) {
		return false, fmt.Errorf("vertex '%s' already labeled by '%s' for correlationIDs '%v'",
			v.ID, r.ResultLabel, v.GetBlock(r.ResultLabel))
	}
	return true, nil
}

type ParentVertexLabelingRule struct {
	Name        string
	CurrentType string
	ParentType  string
	ResultLabel string
	Exclusive   bool
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
			if !visitedVertex.ReserveBlock(r.ResultLabel, cID, r.Exclusive) {
				return fmt.Errorf("vertex '%s' already labeled by '%s' for correlationIDs '%v'",
					visitedVertex.ID, r.ResultLabel, v.GetBlock(r.ResultLabel))
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
	Exclusive   bool
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
			if !visitedVertex.ReserveBlock(r.ResultLabel, cID, r.Exclusive) {
				return fmt.Errorf("vertex '%s' already labeled by '%s' for correlationID '%v'",
					visitedVertex.ID, r.ResultLabel, v.GetBlock(r.ResultLabel))
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
	Exclusive   bool
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
	markChildrenRule := ChildrenVertexLabelingRule{"internal for BrotherVertexLabelingRule",
		r.ParentType, r.BrotherType, r.ResultLabel, r.Exclusive}
	return markChildrenRule.ApplyRule(parent, cID)
}
