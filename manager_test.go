package gographlabel

import (
	"testing"

	"github.com/guntenbein/gojsonut"
)

const Company = "COMPANY"
const UploadChannel = "UPLOAD_CHANNEL"
const Listing = "LISTING"
const Hold = "HOLD"

func initSimpleDefault() (*Vertex, Manager) {
	company := makeCompany()

	ruler := make(Ruler)

	ruler.Add(Default,
		// vertex type assertion
		CurrentVertexTypeCheckingRule{"'default' action types accepted", []string{UploadChannel, Listing}},
		// default
		CurrentVertexLabelingRule{"1.1.1", "", Default, true},
		ParentVertexLabelingRule{"1.1.2", "", UploadChannel, Default, true},
		ChildrenVertexLabelingRule{"1.1.3", UploadChannel, "", Default, true},
		BrotherVertexLabelingRule{"1.1.4", "", UploadChannel, "", Default, true},
		// negotiate
		//CurrentVertexLabelingRule{"1.2.1", Listing, Negotiate},
		//ChildrenVertexLabelingRule{"1.2.2", UploadChannel, Listing, Negotiate},
		//BrotherVertexLabelingRule{"1.2.3", Listing, UploadChannel, Listing, Default},
		// externalAPI
		ParentVertexLabelingRule{"1.3.1", "", Company, ExternalAPI, false},
	)

	ruler.Add(Negotiate,
		// vertex type assertion
		CurrentVertexTypeCheckingRule{"'negotiate' action types accepted", []string{Listing}},
		// default
		// -
		// negotiate
		CurrentVertexLabelingRule{"2.2.1", Listing, Negotiate, true},
		// externalAPI
		ParentVertexLabelingRule{"2.3.1", Listing, Company, ExternalAPI, false},
	)

	ruler.Add(ExternalAPI,
		// vertex type assertion
		CurrentVertexTypeCheckingRule{"'externalAPI' action types accepted", []string{Company}},
		// default
		ChildrenVertexLabelingRule{"3.1.1", Company, UploadChannel, Default, false},
		ChildrenVertexLabelingRule{"3.1.2", Company, Listing, Default, false},
		ChildrenVertexLabelingRule{"3.1.2", Company, Hold, Default, false},
		// negotiate
		ChildrenVertexLabelingRule{"3.2.1", Company, Listing, Negotiate, false},
		// externalAPI
		CurrentVertexLabelingRule{"2.2.1", Company, ExternalAPI, true},
	)

	manager := MakeManager(ruler)
	return company, manager
}

func makeCompany() *Vertex {
	company := NewVertex("COMPANY01", Company)
	uploadChannel := NewVertex("UPLOAD_CHANNEL01", UploadChannel)
	listing1 := NewVertex("LISTING01", Listing)
	listing2 := NewVertex("LISTING02", Listing)
	hold := NewVertex("HOLD01", Hold)

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

func TestActionAppliedToWrongType(t *testing.T) {
	company, manager := initSimpleDefault()

	// default
	if err := manager.CalculateBlocks(company, BlockOrder{Default, "COMPANY01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}
	if err := manager.CalculateBlocks(company, BlockOrder{Default, "HOLD01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}

	// negotiate
	if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "COMPANY01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}
	if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "HOLD01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}
	if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "UPLOAD_CHANNEL01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}

	// external API
	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "LISTING01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}
	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "HOLD01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}
	if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "UPLOAD_CHANNEL01", "process01"}); err == nil {
		t.Fatalf("error expected, but doeas not happen")
	}

}

func BenchmarkApplyingDefaultRules(b *testing.B) {
	b.ReportAllocs()
	company, manager := initSimpleDefault()
	for i := 0; i < b.N; i++ {
		if err := manager.CalculateBlocks(company, BlockOrder{Default, "LISTING01", "process01"}); err != nil {
			b.Fatalf("error happens, but unexpected: %s", err)
		}
	}
}

func BenchmarkApplyingNegotiateRules(b *testing.B) {
	b.ReportAllocs()
	company, manager := initSimpleDefault()
	for i := 0; i < b.N; i++ {
		if err := manager.CalculateBlocks(company, BlockOrder{Negotiate, "LISTING01", "process01"}); err != nil {
			b.Fatalf("error happens, but unexpected: %s", err)
		}
	}
}

func BenchmarkApplyingExternalAPIRules(b *testing.B) {
	b.ReportAllocs()
	company, manager := initSimpleDefault()
	for i := 0; i < b.N; i++ {
		if err := manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "COMPANY01", "process01"}); err != nil {
			b.Fatalf("error happens, but unexpected: %s", err)
		}
	}
}
