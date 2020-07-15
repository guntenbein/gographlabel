package hierarchy_repository

import "github.com/guntenbein/gographlabel"

type HierarchyProvider interface {
	ReadHierarchy(hierarchyId string) (hierarchy *gographlabel.Vertex, err error)
}

type StubHierarchyProvider struct {
	makeHierarchy func() *gographlabel.Vertex
}

func NewStubHierarchyProvider(makeHierarchy func() *gographlabel.Vertex) StubHierarchyProvider {
	return StubHierarchyProvider{makeHierarchy}
}

func (imhp StubHierarchyProvider) ReadHierarchy(_ string) (hierarchy *gographlabel.Vertex, err error) {
	return imhp.makeHierarchy(), nil
}
