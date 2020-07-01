package gographlabel

import (
	"testing"

	"github.com/guntenbein/gojsonut"
)

func initSimpleDefault() (*Vertex, Manager) {
	company := makeCompany()

	ruler := make(Ruler)

	ruler.Add(Default,
		// TODO have a separation by not input action only, but by result label as well
		// default
		CurrentVertexLabelingRule{"1.1.1", "", Default},
		ParentVertexLabelingRule{"1.1.2", "", "UPLOAD_CHANNEL", Default},
		ChildrenVertexLabelingRule{"1.1.3", "UPLOAD_CHANNEL", "", Default},
		BrotherVertexLabelingRule{"1.1.4", "", "UPLOAD_CHANNEL", "", Default},
		// negotiate
		//CurrentVertexLabelingRule{"1.2.1", "LISTING", Negotiate},
		//ChildrenVertexLabelingRule{"1.2.2", "UPLOAD_CHANNEL", "LISTING", Negotiate},
		//BrotherVertexLabelingRule{"1.2.3", "LISTING", "UPLOAD_CHANNEL", "LISTING", Default},
		// externalAPI
		ParentVertexLabelingRule{"1.3.1", "", "COMPANY", ExternalAPI},
		BrotherVertexLabelingRule{"1.3.2", "", "COMPANY", "", ExternalAPI},
	)

	ruler.Add(Negotiate,
		// default
		// -
		// negotiate
		CurrentVertexLabelingRule{"2.2.1", "LISTING", Negotiate},
		// externalAPI
		// -
	)

	ruler.Add(ExternalAPI,
		// default
		ChildrenVertexLabelingRule{"3.1.1", "COMPANY", "UPLOAD_CHANNEL", Default},
		ChildrenVertexLabelingRule{"3.1.2", "COMPANY", "LISTING", Default},
		ChildrenVertexLabelingRule{"3.1.2", "COMPANY", "HOLD", Default},
		// negotiate
		ChildrenVertexLabelingRule{"3.2.1", "COMPANY", "LISTING", Negotiate},
		// externalAPI
		CurrentVertexLabelingRule{"2.2.1", "COMPANY", ExternalAPI},
	)

	manager := MakeManager(ruler)
	return company, manager
}

func makeCompany() *Vertex {
	company := NewVertex("COMPANY01", "COMPANY")
	uploadChannel := NewVertex("UPLOAD_CHANNEL01", "UPLOAD_CHANNEL")
	listing1 := NewVertex("LISTING01", "LISTING")
	listing2 := NewVertex("LISTING02", "LISTING")
	hold := NewVertex("HOLD01", "HOLD")

	// relate them in both directions
	company.AddChildren(uploadChannel.AddChildren(listing1, listing2, hold))
	return company
}

func TestBlockByDefaultAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
    "data": {
        "id": "COMPANY01",
        "type": "COMPANY"
    },
    "labels": {
        "externalAPI": "process01"
    },
    "children": [
        {
            "data": {
                "id": "UPLOAD_CHANNEL01",
                "type": "UPLOAD_CHANNEL"
            },
            "labels": {
                "default": "process01",
                "externalAPI": "process01"
            },
            "children": [
                {
                    "data": {
                        "id": "LISTING01",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": "process01",
                        "externalAPI": "process01"
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "LISTING02",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": "process01",
                        "externalAPI": "process01"
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "HOLD01",
                        "type": "HOLD"
                    },
                    "labels": {
                        "default": "process01",
                        "externalAPI": "process01"
                    },
                    "children": null
                }
            ]
        }
    ]
}`

	if err := manager.CalculateBlocks(company, BlockOrder{Default, "UPLOAD_CHANNEL01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	gojsonut.JsonCompare(t, company, expected, false)

	if err := manager.CalculateBlocks(company, BlockOrder{Default, "LISTING01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	gojsonut.JsonCompare(t, company, expected, false)

}

func TestBlockByNegotiateAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
    "data": {
        "id": "COMPANY01",
        "type": "COMPANY"
    },
    "labels": {},
    "children": [
        {
            "data": {
                "id": "UPLOAD_CHANNEL01",
                "type": "UPLOAD_CHANNEL"
            },
            "labels": {},
            "children": [
                {
                    "data": {
                        "id": "LISTING01",
                        "type": "LISTING"
                    },
                    "labels": {
                        "negotiate": "process01"
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "LISTING02",
                        "type": "LISTING"
                    },
                    "labels": {},
                    "children": null
                },
                {
                    "data": {
                        "id": "HOLD01",
                        "type": "HOLD"
                    },
                    "labels": {},
                    "children": null
                }
            ]
        }
    ]
}`

	if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "LISTING01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	gojsonut.JsonCompare(t, company, expected, false)

}

func TestBlockByExternalAPIAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
    "data": {
        "id": "COMPANY01",
        "type": "COMPANY"
    },
    "labels": {
        "externalAPI": "process01"
    },
    "children": [
        {
            "data": {
                "id": "UPLOAD_CHANNEL01",
                "type": "UPLOAD_CHANNEL"
            },
            "labels": {
                "default": "process01"
            },
            "children": [
                {
                    "data": {
                        "id": "LISTING01",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": "process01",
                        "negotiate": "process01"
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "LISTING02",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": "process01",
                        "negotiate": "process01"
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "HOLD01",
                        "type": "HOLD"
                    },
                    "labels": {
                        "default": "process01"
                    },
                    "children": null
                }
            ]
        }
    ]
}`

	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	gojsonut.JsonCompare(t, company, expected, false)

}

func TestBlockByManyActionsSuccess(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
    "data": {
        "id": "COMPANY01",
        "type": "COMPANY"
    },
    "labels": {
        "externalAPI": "process02"
    },
    "children": [
        {
            "data": {
                "id": "UPLOAD_CHANNEL01",
                "type": "UPLOAD_CHANNEL"
            },
            "labels": {
                "default": "process02",
                "externalAPI": "process02"
            },
            "children": [
                {
                    "data": {
                        "id": "LISTING01",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": "process02",
                        "externalAPI": "process02",
                        "negotiate": "process01"
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "LISTING02",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": "process02",
                        "externalAPI": "process02"
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "HOLD01",
                        "type": "HOLD"
                    },
                    "labels": {
                        "default": "process02",
                        "externalAPI": "process02"
                    },
                    "children": null
                }
            ]
        }
    ]
}`

	if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "LISTING01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	if err := manager.CalculateBlocks(company, BlockOrder{Default, "LISTING01", "process02"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	gojsonut.JsonCompare(t, company, expected, false)

}

func TestBlockByManyActionsFail01(t *testing.T) {
	company, manager := initSimpleDefault()

	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	if err := manager.CalculateBlocks(company, BlockOrder{Default, "LISTING01", "process02"}); err == nil {
		t.Fatalf("error does not happen, but expected")
	}

}

func TestBlockByManyActionsFail02(t *testing.T) {
	company, manager := initSimpleDefault()

	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "LISTING01", "process02"}); err == nil {
		t.Fatalf("error does not happen, but expected")
	}

}
