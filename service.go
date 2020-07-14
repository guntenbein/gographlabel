package gographlabel

import (
	"context"
	"errors"
)

type BlockProvider interface {
	ReadBlocks(hierarchyId string) (blocks []BlockOrder, err error)
	WriteBlocks(hierarchyId string, blocks []BlockOrder) error
}

type HierarchyProvider interface {
	ReadHierarchy(hierarchyId string) (hierarchy *Vertex, err error)
}

type HierarchyLocker interface {
	LockHierarchy(hierarchyId string) (err error)
	UnlockHierarchy(hierarchyId string) (err error)
}

type Service struct {
	blockProvider     BlockProvider
	hierarchyProvider HierarchyProvider
	locker            HierarchyLocker
	manager           Manager
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
	hierarchy, err := s.hierarchyProvider.ReadHierarchy(req.HierarchyId)
	if err != nil {
		return err // wrap
	}
	s.locker.LockHierarchy(req.HierarchyId)
	defer s.locker.UnlockHierarchy(req.HierarchyId)
	blockOrders, err := s.blockProvider.ReadBlocks(req.HierarchyId)
	if err != nil {
		return err // wrap
	}
	arrivedOrder := BlockOrder{req.Action, req.VertexID, req.CorrelationID}
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

func (s Service) Unblock(ctx context.Context, req BlockRequest) error {
	s.locker.LockHierarchy(req.HierarchyId)
	defer s.locker.UnlockHierarchy(req.HierarchyId)
	blockOrders, err := s.blockProvider.ReadBlocks(req.HierarchyId)
	resultingBlockOrders := make([]BlockOrder, 0)
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

func (s Service) IsBlocked(ctx context.Context, req BlockRequest) (bool, error) {
	hierarchy, err := s.hierarchyProvider.ReadHierarchy(req.HierarchyId)
	if err != nil {
		return false, err // wrap
	}
	vertex, err := hierarchy.FindById(req.VertexID)
	if err != nil {
		return false, err // wrap
	}
	if vertex == nil {
		return false, nil
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
		Hierarchy *Vertex
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
