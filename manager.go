package gographlabel

import (
	"context"
	"fmt"
	"time"
)

type Rule interface {
	ApplyRule(toVertex *Vertex) (bool, error)
}

func ApplyRules(v *Vertex, rr ...Rule) error {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	for {
		ruleApplied, err := applyRulesOnce(v, rr...)
		if err != nil {
			return err
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
	return nil
}

func applyRulesOnce(v *Vertex, rr ...Rule) (bool, error) {
	ruleApplied := false
	for _, r := range rr {
		applied, err := r.ApplyRule(v)
		if err != nil {
			return false, err
		}
		ruleApplied = ruleApplied || applied
	}
	for _, cv := range v.Children {
		applied, err := applyRulesOnce(cv, rr...)
		if err != nil {
			return false, err
		}
		ruleApplied = ruleApplied || applied
	}
	return ruleApplied, nil
}
