package gographlabel

type LabelBrokerEnum map[string]*Block

type Block struct {
	CorrelationIds map[string]struct{}
	Exclusive      bool
}

func (le LabelBrokerEnum) GetBlock(label string) *Block {
	block, ok := le[label]
	if ok {
		return block
	}
	emptyBlock := Block{make(map[string]struct{}), false}
	le[label] = &emptyBlock
	return &emptyBlock
}

func (le LabelBrokerEnum) ReserveBlock(label, correlationID string, exclusive bool) bool {
	block := le.GetBlock(label)
	_, ok := block.CorrelationIds[correlationID]
	if ok {
		return true
	}
	if block.Exclusive {
		return false
	}
	if exclusive && len(block.CorrelationIds) > 0 {
		return false
	}
	block.CorrelationIds[correlationID] = struct{}{}
	block.Exclusive = exclusive
	return true
}
