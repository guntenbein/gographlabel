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
		CurrentVertexLabelingRule{"1.1.1", "", Default, true},
		ParentVertexLabelingRule{"1.1.2", "", "UPLOAD_CHANNEL", Default, true},
		ChildrenVertexLabelingRule{"1.1.3", "UPLOAD_CHANNEL", "", Default, true},
		BrotherVertexLabelingRule{"1.1.4", "", "UPLOAD_CHANNEL", "", Default, true},
		// negotiate
		//CurrentVertexLabelingRule{"1.2.1", "LISTING", Negotiate},
		//ChildrenVertexLabelingRule{"1.2.2", "UPLOAD_CHANNEL", "LISTING", Negotiate},
		//BrotherVertexLabelingRule{"1.2.3", "LISTING", "UPLOAD_CHANNEL", "LISTING", Default},
		// externalAPI
		ParentVertexLabelingRule{"1.3.1", "", "COMPANY", ExternalAPI, false},
	)

	ruler.Add(Negotiate,
		// default
		// -
		// negotiate
		CurrentVertexLabelingRule{"2.2.1", "LISTING", Negotiate, true},
		// externalAPI
		ParentVertexLabelingRule{"2.3.1", "LISTING", "COMPANY", ExternalAPI, false},
	)

	ruler.Add(ExternalAPI,
		// default
		ChildrenVertexLabelingRule{"3.1.1", "COMPANY", "UPLOAD_CHANNEL", Default, false},
		ChildrenVertexLabelingRule{"3.1.2", "COMPANY", "LISTING", Default, false},
		ChildrenVertexLabelingRule{"3.1.2", "COMPANY", "HOLD", Default, false},
		// negotiate
		ChildrenVertexLabelingRule{"3.2.1", "COMPANY", "LISTING", Negotiate, false},
		// externalAPI
		CurrentVertexLabelingRule{"2.2.1", "COMPANY", ExternalAPI, true},
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
        "externalAPI": {
            "CorrelationIds": {
                "process01": {}
            },
            "Exclusive": false
        }
    },
    "children": [
        {
            "data": {
                "id": "UPLOAD_CHANNEL01",
                "type": "UPLOAD_CHANNEL"
            },
            "labels": {
                "default": {
                    "CorrelationIds": {
                        "process01": {}
                    },
                    "Exclusive": true
                }
            },
            "children": [
                {
                    "data": {
                        "id": "LISTING01",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": true
                        }
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "LISTING02",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": true
                        }
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "HOLD01",
                        "type": "HOLD"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": true
                        }
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

	gojsonut.JsonCompare(t, company, expected)

	if err := manager.CalculateBlocks(company, BlockOrder{Default, "LISTING01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	gojsonut.JsonCompare(t, company, expected)

}

func TestBlockByNegotiateAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
    "data": {
        "id": "COMPANY01",
        "type": "COMPANY"
    },
    "labels": {
        "externalAPI": {
            "CorrelationIds": {
                "process01": {}
            },
            "Exclusive": false
        }
    },
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
                        "negotiate": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": true
                        }
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

	gojsonut.JsonCompare(t, company, expected)

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
        "externalAPI": {
            "CorrelationIds": {
                "process01": {}
            },
            "Exclusive": true
        }
    },
    "children": [
        {
            "data": {
                "id": "UPLOAD_CHANNEL01",
                "type": "UPLOAD_CHANNEL"
            },
            "labels": {
                "default": {
                    "CorrelationIds": {
                        "process01": {}
                    },
                    "Exclusive": false
                }
            },
            "children": [
                {
                    "data": {
                        "id": "LISTING01",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": false
                        },
                        "negotiate": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": false
                        }
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "LISTING02",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": false
                        },
                        "negotiate": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": false
                        }
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "HOLD01",
                        "type": "HOLD"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": false
                        }
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

	gojsonut.JsonCompare(t, company, expected)

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
        "externalAPI": {
            "CorrelationIds": {
                "process01": {},
                "process02": {}
            },
            "Exclusive": false
        }
    },
    "children": [
        {
            "data": {
                "id": "UPLOAD_CHANNEL01",
                "type": "UPLOAD_CHANNEL"
            },
            "labels": {
                "default": {
                    "CorrelationIds": {
                        "process02": {}
                    },
                    "Exclusive": true
                }
            },
            "children": [
                {
                    "data": {
                        "id": "LISTING01",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process02": {}
                            },
                            "Exclusive": true
                        },
                        "negotiate": {
                            "CorrelationIds": {
                                "process01": {}
                            },
                            "Exclusive": true
                        }
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "LISTING02",
                        "type": "LISTING"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process02": {}
                            },
                            "Exclusive": true
                        }
                    },
                    "children": null
                },
                {
                    "data": {
                        "id": "HOLD01",
                        "type": "HOLD"
                    },
                    "labels": {
                        "default": {
                            "CorrelationIds": {
                                "process02": {}
                            },
                            "Exclusive": true
                        }
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

	gojsonut.JsonCompare(t, company, expected)

}

func TestBlockByIncompatibleActionsFail01(t *testing.T) {
	company, manager := initSimpleDefault()

	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	if err := manager.CalculateBlocks(company, BlockOrder{Default, "LISTING01", "process02"}); err == nil {
		t.Fatalf("error does not happen, but expected")
	}

}

func TestBlockByIncompatibleActionsFail02(t *testing.T) {
	company, manager := initSimpleDefault()

	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "LISTING01", "process02"}); err == nil {
		t.Fatalf("error does not happen, but expected")
	}

}

func TestBlockTwiceSameAction(t *testing.T) {
	company, manager := initSimpleDefault()

	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
		t.Fatalf("error happens, but not expected: %v", err)
	}

}
