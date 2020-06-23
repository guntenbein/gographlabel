package gographlabel

import (
	"context"
	"fmt"
	"time"
)

type Vertex struct {
	Type     string
	Parent   *Vertex
	Children []Vertex
	Labels   []string
}

type Rule struct {
	Name           string
	Up             bool
	Recursive      bool
	ConditionType  string
	ConditionLabel string
	ResultLabel    string
}

func (v *Vertex) Contains(label string) bool {
	for _, l := range v.Labels {
		if l == label {
			return true
		}
	}
	return false
}

func (v *Vertex) ApplyRules(rr ...Rule) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	for {
		ruleApplied := false
		for _, r := range rr {
			ruleApplied = ruleApplied || v.applyRule(r)
		}
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

func (v *Vertex) applyRule(r Rule) bool {
	if r.Up {
		if r.Recursive {
			for p := v.Parent; p != nil; {
				if v.applyParent(r, p) {
					return true
				}
			}
		}
		return v.applyParent(r, v.Parent)
	}
	// TODO down part
	return true
}

func (v *Vertex) applyParent(r Rule, parent *Vertex) bool {
	if parent.Contains(r.ConditionLabel) && parent.Type == r.ConditionType {
		v.Labels = append(v.Labels, r.ResultLabel)
		fmt.Printf("applied rule %s for the vertex %+v because of vertext %+v \n", r.Name, v, parent)
		return true
	}
	return false
}
