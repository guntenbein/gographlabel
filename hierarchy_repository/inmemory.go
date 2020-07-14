package hierarchy_repository

import "github.com/guntenbein/gographlabel"

type HierarchyProvider interface {
	ReadHierarchy(hierarchyId string) (hierarchy *gographlabel.Vertex, err error)
}

type InMemoryHierarchyProvider struct {
	hierarchy *gographlabel.Vertex
}

func NewInMemoryHierarchyProvider(hierarchy *gographlabel.Vertex) InMemoryHierarchyProvider {
	return InMemoryHierarchyProvider{hierarchy}
}

func (imhp InMemoryHierarchyProvider) ReadHierarchy(hierarchyId string) (hierarchy *gographlabel.Vertex, err error) {
	return imhp.hierarchy, nil
}
