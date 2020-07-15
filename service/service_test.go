package service

import (
	"context"
	"testing"

	"github.com/guntenbein/gographlabel"
	"github.com/guntenbein/gographlabel/block_repository"
	"github.com/guntenbein/gographlabel/hierarchy_repository"
	"github.com/guntenbein/gographlabel/hierarchylock"
)

func TestServiceBlockCheckSuccess(t *testing.T) {
	service := initServiceTestData()

	blockRequest := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 1",
	}

	err := service.Block(context.Background(), blockRequest)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

	rightBlock, err := service.Check(context.Background(), blockRequest)
	if !rightBlock {
		t.Fatalf("the check for the block should return ok, but the block has not passed the check")
	}
	if err != nil {
		t.Fatalf("should be error, but: %s", err.Error())
	}

}

func TestServiceBlockCheckFailure(t *testing.T) {
	service := initServiceTestData()

	blockRequestProcess1 := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 1",
	}

	err := service.Block(context.Background(), blockRequestProcess1)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

	blockRequestProcess2 := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 2",
	}

	rightBlock, err := service.Check(context.Background(), blockRequestProcess2)
	if rightBlock {
		t.Fatalf("the check for the block should return false, but the block has passed the check")
	}
	if err != nil {
		t.Fatalf("should be error, but: %s", err.Error())
	}

}

func TestServiceBlockCheckListing(t *testing.T) {
	service := initServiceTestData()

	blockRequestProcess1 := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 1",
	}

	err := service.Block(context.Background(), blockRequestProcess1)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

	blockRequestProcess1Listing := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "LISTING01",
		CorrelationID: "process 1",
	}

	rightBlock, err := service.Check(context.Background(), blockRequestProcess1Listing)
	if !rightBlock {
		t.Fatalf("the check for the block should return ok, but the block has not passed the check")
	}
	if err != nil {
		t.Fatalf("should be error, but: %s", err.Error())
	}

	blockRequestProcess2Listing := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "LISTING01",
		CorrelationID: "process 2",
	}

	rightBlock, err = service.Check(context.Background(), blockRequestProcess2Listing)
	if rightBlock {
		t.Fatalf("the check for the block should return false, but the block has passed the check")
	}
	if err != nil {
		t.Fatalf("should be error, but: %s", err.Error())
	}

}

func TestServiceBlockBlockSuccess(t *testing.T) {
	service := initServiceTestData()

	blockRequest := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 1",
	}

	err := service.Block(context.Background(), blockRequest)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

	err = service.Block(context.Background(), blockRequest)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

}

func TestServiceBlockBlockFailure(t *testing.T) {
	service := initServiceTestData()

	blockRequestProcess1 := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 1",
	}

	err := service.Block(context.Background(), blockRequestProcess1)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

	blockRequestProcess2 := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 2",
	}

	err = service.Block(context.Background(), blockRequestProcess2)
	if err == nil {
		t.Fatalf("error not happened, but expected")
	}

}

func TestServiceBlockReleaseBlockSuccess(t *testing.T) {
	service := initServiceTestData()

	blockRequestProcess1 := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 1",
	}

	err := service.Block(context.Background(), blockRequestProcess1)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

	err = service.Unblock(context.Background(), blockRequestProcess1)
	if err != nil {
		t.Fatalf("error happened, but not expected: %s", err.Error())
	}

	blockRequestProcess2 := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "UPLOAD_CHANNEL01",
		CorrelationID: "process 2",
	}

	err = service.Block(context.Background(), blockRequestProcess2)
	if err != nil {
		t.Fatalf("error not expected, but happened: %s", err.Error())
	}

}

func initServiceTestData() Service {
	manager := gographlabel.MakeManager(gographlabel.DefaultRuler())

	makeCompanyHierarchy := func() *gographlabel.Vertex {

		company := gographlabel.NewVertex("COMPANY01", gographlabel.Company)
		uploadChannel := gographlabel.NewVertex("UPLOAD_CHANNEL01", gographlabel.UploadChannel)
		listing1 := gographlabel.NewVertex("LISTING01", gographlabel.Listing)
		listing2 := gographlabel.NewVertex("LISTING02", gographlabel.Listing)
		hold := gographlabel.NewVertex("HOLD01", gographlabel.Hold)

		// relate them in both directions
		company.AddChildren(uploadChannel.AddChildren(listing1, listing2, hold))

		return company
	}

	service := MakeService(
		block_repository.NewInMemoryBlockRepository(),
		hierarchy_repository.NewStubHierarchyProvider(makeCompanyHierarchy),
		hierarchylock.NewSingleServiceLock(),
		manager)
	return service
}

func TestServiceNotFoundError(t *testing.T) {
	service := initServiceTestData()

	blockRequest := BlockRequest{
		HierarchyId:   "COMPANY01",
		Action:        gographlabel.Default,
		VertexID:      "SOMETHING",
		CorrelationID: "process 1",
	}

	err := service.Block(context.Background(), blockRequest)
	if err == nil {
		t.Fatalf("error expected, but not happened:")
	}

	_, err = service.Check(context.Background(), blockRequest)
	if err == nil {
		t.Fatalf("error expected, but not happened:")
	}

	err = service.Unblock(context.Background(), blockRequest)
	if err == nil {
		t.Fatalf("error expected, but not happened:")
	}

}
