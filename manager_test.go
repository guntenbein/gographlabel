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
		CurrentVertexLabelingRule{"1.2.1", "LISTING", Negotiate},
		ChildrenVertexLabelingRule{"1.2.2", "UPLOAD_CHANNEL", "LISTING", Negotiate},
		BrotherVertexLabelingRule{"1.2.3", "LISTING", "UPLOAD_CHANNEL", "LISTING", Default},
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
	company := NewVertex("1", "COMPANY")
	uploadChannel := NewVertex("2", "UPLOAD_CHANNEL")
	listing1 := NewVertex("3", "LISTING")
	listing2 := NewVertex("4", "LISTING")
	hold := NewVertex("5", "HOLD")

	// relate them in both directions
	company.AddChildren(uploadChannel.AddChildren(listing1, listing2, hold))
	return company
}

func TestBlockByDefaultAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
  "data": {
    "id": "1",
    "type": "COMPANY"
  },
  "labels": {
    "externalAPI": {}
  },
  "children": [
    {
      "data": {
        "id": "2",
        "type": "UPLOAD_CHANNEL"
      },
      "labels": {
        "default": {},
        "externalAPI": {}
      },
      "children": [
        {
          "data": {
            "id": "3",
            "type": "LISTING"
          },
          "labels": {
            "default": {},
            "externalAPI": {},
            "negotiate": {}
          },
          "children": null
        },
        {
          "data": {
            "id": "4",
            "type": "LISTING"
          },
          "labels": {
            "default": {},
            "externalAPI": {},
            "negotiate": {}
          },
          "children": null
        },
        {
          "data": {
            "id": "5",
            "type": "HOLD"
          },
          "labels": {
            "default": {},
            "externalAPI": {}
          },
          "children": null
        }
      ]
    }
  ]
}`

	manager.CalculateBlocks(company, BlockOrder{Default, "2"})

	gojsonut.JsonCompare(t, company, expected, false)

	manager.CalculateBlocks(company, BlockOrder{Default, "3"})

	gojsonut.JsonCompare(t, company, expected, false)

}

func TestBlockByNegotiateAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
  "data": {
    "id": "1",
    "type": "COMPANY"
  },
  "labels": {},
  "children": [
    {
      "data": {
        "id": "2",
        "type": "UPLOAD_CHANNEL"
      },
      "labels": {},
      "children": [
        {
          "data": {
            "id": "3",
            "type": "LISTING"
          },
          "labels": {
            "negotiate": {}
          },
          "children": null
        },
        {
          "data": {
            "id": "4",
            "type": "LISTING"
          },
          "labels": {},
          "children": null
        },
        {
          "data": {
            "id": "5",
            "type": "HOLD"
          },
          "labels": {},
          "children": null
        }
      ]
    }
  ]
}`

	manager.CalculateBlocks(company, BlockOrder{Negotiate, "3"})

	gojsonut.JsonCompare(t, company, expected, false)

}

func TestBlockByExternalAPIAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{
  "data": {
    "id": "1",
    "type": "COMPANY"
  },
  "labels": {
    "externalAPI": {}
  },
  "children": [
    {
      "data": {
        "id": "2",
        "type": "UPLOAD_CHANNEL"
      },
      "labels": {
        "default": {}
      },
      "children": [
        {
          "data": {
            "id": "3",
            "type": "LISTING"
          },
          "labels": {
            "default": {},
            "negotiate": {}
          },
          "children": null
        },
        {
          "data": {
            "id": "4",
            "type": "LISTING"
          },
          "labels": {
            "default": {},
            "negotiate": {}
          },
          "children": null
        },
        {
          "data": {
            "id": "5",
            "type": "HOLD"
          },
          "labels": {
            "default": {}
          },
          "children": null
        }
      ]
    }
  ]
}`

	manager.CalculateBlocks(company, BlockOrder{ExternalAPI, "1"})

	gojsonut.JsonCompare(t, company, expected, false)

}
