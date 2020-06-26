package gographlabel

import (
	"errors"
)

type ParentRule struct {
	Name         string
	Recursive    bool
	CurrentType  string
	CurrentLabel string
	ParentType   string
	ParentLabel  string
	ResultLabel  string
}

var LOOP_IN_HIERARCHY_ERROR = errors.New("loops are not allowed in hierarchy")

func (pr *ParentRule) ApplyRule(v *Vertex) (bool, error) {
	if pr.CurrentType != "" && pr.CurrentType != v.Type {
		return false, nil
	}
	if pr.CurrentLabel != "" && !v.Labels.Contains(Label{pr.CurrentLabel}) {
		return false, nil
	}
	exploredVertexes := make(map[*Vertex]struct{})
	exploredVertexes[v] = struct{}{}
	for p := v.Parent; p != nil; {
		if _, ok := exploredVertexes[p]; !ok {
			return false, LOOP_IN_HIERARCHY_ERROR
		}
		exploredVertexes[p] = struct{}{}
		if checkTypeLabelAndApplyLabel(v, p, pr.ParentLabel, pr.ParentType, pr.ResultLabel) {
			return true, nil
		}
		if !pr.Recursive {
			return false, nil
		}
		p = p.Parent
	}
	return false, nil
}

type ChildRule struct {
	Name         string
	Recursive    bool
	CurrentType  string
	CurrentLabel string
	ChildType    string
	ChildLabel   string
	ResultLabel  string
}

func (chr *ChildRule) ApplyRule(v *Vertex) (bool, error) {
	if chr.CurrentType != "" && chr.CurrentType != v.Type {
		return false, nil
	}
	if chr.CurrentLabel != "" && !v.Labels.Contains(Label{chr.CurrentLabel}) {
		return false, nil
	}

	if !chr.Recursive {
		for i, _ := range v.Children {
			if checkTypeLabelAndApplyLabel(v, v.Children[i], chr.CurrentLabel, chr.ChildType, chr.ResultLabel) {
				return true, nil
			}
		}
	}
	// TODO comparison operations for Vertex !!!!
	exploredVertexes := make(map[VertexData]struct{})
	exploredVertexes[v.VertexData] = struct{}{}

	queue := []*Vertex{}
	queue = append(queue, v)
	return chr.bfs(v, queue, exploredVertexes)
}

// bfs - recursive function for implementing bfs for our hierarchy
// TODO - workers
func (chr *ChildRule) bfs(initialV *Vertex, queue []*Vertex, exploredVertexes map[VertexData]struct{}) (bool, error) {
	//This appends the value of the root of subtree or tree to the result
	if len(queue) == 0 {
		return false, nil
	}
	currentV := queue[0]
	if _, ok := exploredVertexes[currentV.VertexData]; ok {
		return false, LOOP_IN_HIERARCHY_ERROR
	}
	if checkTypeLabelAndApplyLabel(initialV, currentV, chr.ChildLabel, chr.ChildType, chr.ResultLabel) {
		return true, nil
	}
	if len(currentV.Children) > 0 {
		queue = append(queue, currentV.Children...)
	}
	return chr.bfs(initialV, queue[1:], exploredVertexes)
}

func checkTypeLabelAndApplyLabel(v *Vertex, r *Vertex, rLabel, rType, tobeLabel string) bool {
	if r.Contains(rLabel) && r.Type == rType {
		v.Labels[Label{tobeLabel}] = struct{}{}
		return true
	}
	return false
}
