package gographlabel

import (
	"testing"
)

func TestFindById(t *testing.T) {
	company := NewVertex("1", "COMPANY")
	uploadChannel := NewVertex("2", "UPLOAD_CHANNEL")
	listing1 := NewVertex("3", "LISTING")
	listing2 := NewVertex("4", "LISTING")
	hold := NewVertex("5", "HOLD")

	// relate them
	company.AddChildren(uploadChannel.AddChildren(listing1, listing2, hold))

	foundVertex, err := company.FindById("1")
	if err != nil || foundVertex == nil || foundVertex.ID != "1" || foundVertex.Type != "COMPANY" {
		t.Fatalf("find function works unproperly, found %+v", foundVertex)
	}

	foundVertex, err = company.FindById("3")
	if err != nil || foundVertex == nil || foundVertex.ID != "3" || foundVertex.Type != "LISTING" {
		t.Fatalf("find function works unproperly, found %+v", foundVertex)
	}

	foundVertex, err = company.FindById("5")
	if err != nil || foundVertex == nil || foundVertex.ID != "5" || foundVertex.Type != "HOLD" {
		t.Fatalf("find function works unproperly, found %+v", foundVertex)
	}
}
