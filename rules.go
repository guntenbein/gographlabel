package gographlabel

import (
	"context"
	"errors"
	"fmt"
	"time"
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

func (pr *ParentRule) ApplyRule(v *Vertex) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
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
			return false, errors.New("loops are not allowed in hierarchy")
		}
		exploredVertexes[p] = struct{}{}
		if checkTypeLabelAndApplyLabel(v, p, pr.ParentLabel, pr.ParentType, pr.ResultLabel) {
			return true, nil
		}
		if !pr.Recursive {
			return false, nil
		}
		p = p.Parent
		select {
		case <-ctx.Done():
			fmt.Println("exit because of timeout - probable loop for the rule application")
			break
		default:
		}
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
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	if chr.CurrentType != "" && chr.CurrentType != v.Type {
		return false, nil
	}
	if chr.CurrentLabel != "" && !v.Labels.Contains(Label{chr.CurrentLabel}) {
		return false, nil
	}
	exploredVertexes := make(map[*Vertex]struct{})
	exploredVertexes[v] = struct{}{}
	// TODO complete it
	for p := v.Parent; p != nil; {
		if _, ok := exploredVertexes[p]; !ok {
			return false, errors.New("loops are not allowed in hierarchy")
		}
		exploredVertexes[p] = struct{}{}
		if checkTypeLabelAndApplyLabel(v, p, chr.ChildLabel, chr.ChildType, chr.ResultLabel) {
			return true, nil
		}
		if !chr.Recursive {
			return false, nil
		}
		p = p.Parent
		select {
		case <-ctx.Done():
			fmt.Println("exit because of timeout - probable loop for the rule application")
			break
		default:
		}
	}
	return false, nil
}

func checkTypeLabelAndApplyLabel(v *Vertex, r *Vertex, rLabel, rType, tobeLabel string) bool {
	if r.Contains(rLabel) && r.Type == rType {
		v.Labels[Label{tobeLabel}] = struct{}{}
		return true
	}
	return false
}
