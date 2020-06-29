package gographlabel

type Rule interface {
	ApplyRule(toVertex *Vertex) (bool, error)
}

type Ruler map[string][]Rule

func (rr Ruler) Add(action string, rules ...Rule) {
	if len(rules) == 0 {
		return
	}
	storedRules, ok := rr[action]
	if !ok {
		storedRules = make([]Rule, 0)
	}
	storedRules = append(storedRules, rules...)
	rr[action] = storedRules
}

type Manager struct {
	ruler Ruler
}

func MakeManager(ruler Ruler) Manager {
	return Manager{ruler}
}

func (m Manager) CalculateBlocks(hierarchy *Vertex, orders ...BlockOrder) error {
	// todo provide in more functional way - copy the hierarchy and output the changed copied version
	for _, o := range orders {
		rules := m.ruler[o.Action]
		if rules == nil || len(rules) == 0 {
			continue
		}
		orderedVertex, err := hierarchy.FindById(o.ID)
		if err != nil {
			return err
		}
		if orderedVertex == nil {
			continue
		}
		for _, r := range rules {
			r.ApplyRule(orderedVertex)
		}
	}
	return nil
}

type BlockOrder struct {
	Action string
	ID     string
}
