package gographlabel

import "errors"

type Vertex struct {
	VertexData   `json:"data"`
	LabelStorage `json:"labels"`
	Parent       *Vertex   `json:"-"`
	Children     []*Vertex `json:"children"`
}

func NewVertex(id, tp string) *Vertex {
	return &Vertex{VertexData: VertexData{id, tp}, LabelStorage: make(LabelEnum)}
}

func (v *Vertex) AddChildren(children ...*Vertex) *Vertex {
	if children == nil || len(children) == 0 {
		return v
	}
	for i, _ := range children {
		children[i].Parent = v
	}
	if v.Children == nil {
		v.Children = children
		return v
	}
	v.Children = append(v.Children, children...)
	return v
}

func (v *Vertex) FindById(id string) (*Vertex, error) {
	var foundVertex *Vertex
	if err := v.ApplyChildren(func(visitedVertex *Vertex) error {
		if visitedVertex.ID == id {
			foundVertex = visitedVertex
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return foundVertex, nil
}

func (v *Vertex) ApplyChildren(apply func(visitedVertex *Vertex) error) error {
	exploredVertexes := make(map[VertexData]struct{})
	queue := []*Vertex{v}
	return bfs(true, v, queue, exploredVertexes, apply)
}

func (v *Vertex) ApplyParents(apply func(visitedVertex *Vertex) error) error {
	exploredVertexes := make(map[VertexData]struct{})
	queue := []*Vertex{v}
	return bfs(false, v, queue, exploredVertexes, apply)
}

// bfs - recursive function for implementing bfs for our hierarchy
// TODO - workers?
func bfs(down bool, initialVertex *Vertex, queue []*Vertex, exploredVertexes map[VertexData]struct{}, apply func(visitedVertex *Vertex) error) error {
	//This appends the value of the root of subtree or tree to the result
	if len(queue) == 0 {
		return nil
	}
	currentV := queue[0]
	if _, ok := exploredVertexes[currentV.VertexData]; ok {
		return errors.New(LoopInHierarchyError)
	}
	exploredVertexes[currentV.VertexData] = struct{}{}
	if err := apply(currentV); err != nil {
		return err
	}
	if down {
		if len(currentV.Children) > 0 {
			queue = append(queue, currentV.Children...)
		}
	} else {
		if currentV.Parent != nil {
			queue = append(queue, currentV.Parent)
		}
	}
	return bfs(down, initialVertex, queue[1:], exploredVertexes, apply)
}

type VertexData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type LabelStorage interface {
	Get(l string) (string, bool)
	MustGet(l string) string
	Reserve(l, forCorrelationID string) bool
}
