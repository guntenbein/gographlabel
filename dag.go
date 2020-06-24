package gographlabel

import (
	"context"
	"fmt"
	"time"
)

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

type Rule interface {
	ApplyRule(toVertex *Vertex) bool
}

type ParentRule struct {
	Name         string
	Recursive    bool
	CurrentType  string
	CurrentLabel string
	ParentType   string
	ParentLabel  string
	ResultLabel  string
}

func (pr *ParentRule) ApplyRule(v *Vertex) bool {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	if pr.CurrentType != "" && pr.CurrentType != v.Type {
		return false
	}
	if pr.CurrentLabel != "" && !v.Labels.Contains(Label{pr.CurrentLabel}) {
		return false
	}
	for p := v.Parent; p != nil; {
		if p.Contains(pr.ParentLabel) && p.Type == pr.ParentType {
			v.Labels[Label{pr.ResultLabel}] = struct{}{}
			return true
		}
		if !pr.Recursive {
			return false
		}
		p = p.Parent
		select {
		case <-ctx.Done():
			fmt.Println("exit because of timeout - probable loop for the rule application")
			break
		default:
		}
	}
	return false
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

func (chr *ChildRule) ApplyRule(v *Vertex) bool {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	if chr.CurrentType != "" && chr.CurrentType != v.Type {
		return false
	}
	if chr.CurrentLabel != "" && !v.Labels.Contains(Label{chr.CurrentLabel}) {
		return false
	}
	// TODO change it to be really recursive
	for p := v.Parent; p != nil; {
		if p.Contains(chr.ChildLabel) && p.Type == chr.ChildType {
			v.Labels[Label{chr.ResultLabel}] = struct{}{}
			return true
		}
		if !chr.Recursive {
			return false
		}
		p = p.Parent
		select {
		case <-ctx.Done():
			fmt.Println("exit because of timeout - probable loop for the rule application")
			break
		default:
		}
	}
	return false
}

func (v *Vertex) Contains(label string) bool {
	v.Labels.Contains(Label{label})
	return false
}

func (v *Vertex) ApplyRules(rr ...Rule) {
	// TODO to think about timeout more
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	for {
		ruleApplied := v.applyRulesOnce(rr...)
		if !ruleApplied {
			break
		}
		select {
		case <-ctx.Done():
			fmt.Println("exit because of timeout - probable loop for the rules application")
			break
		default:
		}
	}
}

func (v *Vertex) applyRulesOnce(rr ...Rule) bool {
	ruleApplied := false
	for _, r := range rr {
		ruleApplied = ruleApplied || r.ApplyRule(v)
	}
	for _, cv := range v.Children {
		ruleApplied = ruleApplied || cv.applyRulesOnce()
	}
	return ruleApplied
}
