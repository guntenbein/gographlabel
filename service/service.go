package service

import (
	"context"
	"errors"

	"github.com/guntenbein/gographlabel"
)

type BlockProvider interface {
	ReadBlocks(hierarchyId string) (blocks []gographlabel.BlockOrder, err error)
	WriteBlocks(hierarchyId string, blocks []gographlabel.BlockOrder) error
}

type HierarchyProvider interface {
	ReadHierarchy(hierarchyId string) (hierarchy *gographlabel.Vertex, err error)
}

type HierarchyLocker interface {
	LockHierarchy(hierarchyId string) (err error)
	UnlockHierarchy(hierarchyId string) (err error)
}

type Service struct {
	blockProvider     BlockProvider
	hierarchyProvider HierarchyProvider
	locker            HierarchyLocker
	manager           gographlabel.Manager
}

// MakeService - service constructor
func MakeService(blockProvider BlockProvider,
	hierarchyProvider HierarchyProvider,
	locker HierarchyLocker,
	manager gographlabel.Manager,
) Service {
	return Service{
		blockProvider:     blockProvider,
		hierarchyProvider: hierarchyProvider,
		locker:            locker,
		manager:           manager,
	}
}

type (
	BlockRequest struct {
		HierarchyId   string
		Action        string
		VertexID      string
		CorrelationID string
	}
)

func (s Service) Block(ctx context.Context, req BlockRequest) error {
	hierarchy, _, err := s.getHierarchyVertexValidate(req)
	if err != nil {
		return err
	}
	s.locker.LockHierarchy(req.HierarchyId)
	defer s.locker.UnlockHierarchy(req.HierarchyId)
	blockOrders, err := s.blockProvider.ReadBlocks(req.HierarchyId)
	if err != nil {
		return err // wrap
	}
	arrivedOrder := gographlabel.BlockOrder{req.Action, req.VertexID, req.CorrelationID}
	resultingBlockOrders := append(blockOrders, arrivedOrder)
	err = s.manager.CalculateBlocks(hierarchy, resultingBlockOrders...)
	if err != nil {
		return errors.New("error blocking resource") // standard error locked 423
	}
	err = s.blockProvider.WriteBlocks(req.HierarchyId, resultingBlockOrders)
	if err != nil {
		return err // wrap
	}
	return nil
}

func (s Service) getHierarchyVertexValidate(req BlockRequest) (hierarchy *gographlabel.Vertex, vertex *gographlabel.Vertex, err error) {
	hierarchy, err = s.hierarchyProvider.ReadHierarchy(req.HierarchyId)
	if err != nil {
		return nil, nil, err // wrap
	}
	vertex, err = hierarchy.FindById(req.VertexID)
	if err != nil {
		return nil, nil, err // wrap
	}
	if vertex == nil {
		return nil, nil, errors.New("entity not found") // standard entity not found error
	}
	return
}

func (s Service) Unblock(ctx context.Context, req BlockRequest) error {
	_, _, err := s.getHierarchyVertexValidate(req)
	if err != nil {
		return err
	}
	s.locker.LockHierarchy(req.HierarchyId)
	defer s.locker.UnlockHierarchy(req.HierarchyId)
	blockOrders, err := s.blockProvider.ReadBlocks(req.HierarchyId)
	resultingBlockOrders := make([]gographlabel.BlockOrder, 0)
	for _, order := range blockOrders {
		if order.CorrelationID != req.CorrelationID &&
			order.Action != req.Action &&
			order.VertexID != req.VertexID {
			resultingBlockOrders = append(resultingBlockOrders, order)
		}
	}
	if len(resultingBlockOrders) < len(blockOrders) {
		err = s.blockProvider.WriteBlocks(req.HierarchyId, resultingBlockOrders)
		if err != nil {
			return err // wrap
		}
	}
	return nil
}

func (s Service) Check(ctx context.Context, req BlockRequest) (bool, error) {
	hierarchy, vertex, err := s.getHierarchyVertexValidate(req)
	if err != nil {
		return false, err
	}
	s.locker.LockHierarchy(req.HierarchyId)
	defer s.locker.UnlockHierarchy(req.HierarchyId)
	blockOrders, err := s.blockProvider.ReadBlocks(req.HierarchyId)
	if err != nil {
		return false, err // wrap
	}
	err = s.manager.CalculateBlocks(hierarchy, blockOrders...)
	if err != nil {
		return false, err
	}
	block := vertex.GetBlock(req.Action)
	if block == nil || block.CorrelationIds == nil {
		return false, nil
	}
	if _, ok := block.CorrelationIds[req.CorrelationID]; ok {
		return true, nil
	}
	return false, nil
}

type (
	StatusRequest struct {
		HierarchyId string
	}
	StatusResponse struct {
		Hierarchy *gographlabel.Vertex
	}
)

func (s Service) Status(ctx context.Context, req StatusRequest) (StatusResponse, error) {
	hierarchy, err := s.hierarchyProvider.ReadHierarchy(req.HierarchyId)
	if err != nil {
		return StatusResponse{}, err // wrap
	}
	s.locker.LockHierarchy(req.HierarchyId)
	defer s.locker.UnlockHierarchy(req.HierarchyId)
	blockOrders, err := s.blockProvider.ReadBlocks(req.HierarchyId)
	if err != nil {
		return StatusResponse{}, err // wrap
	}
	err = s.manager.CalculateBlocks(hierarchy, blockOrders...)
	if err != nil {
		return StatusResponse{}, err
	}
	return StatusResponse{hierarchy}, nil
}
