package gographlabel

import (
	"testing"

	"github.com/guntenbein/gojsonut"
)

func TestBlockUploadChannelByDefaultAction(t *testing.T) {
	company, manager := initSimpleDefault()

	expected := `
{  
"node": {
    "id": "1",
    "type": "COMPANY"
  },
  "labels": {},
  "children": [
    {
      "node": {
        "id": "2",
        "type": "UPLOAD_CHANNEL"
      },
      "labels": {
        "default": {}
      },
      "children": [
        {
          "node": {
            "id": "3",
            "type": "LISTING"
          },
          "labels": {
            "default": {}
          },
          "children": null
        },
        {
          "node": {
            "id": "4",
            "type": "LISTING"
          },
          "labels": {
            "default": {}
          },
          "children": null
        },
        {
          "node": {
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
}
`

	manager.CalculateBlocks(company, BlockOrder{Default, "2"})

	gojsonut.JsonCompare(t, company, expected, false)

}

func TestBlockListingByDefaultAction(t *testing.T) {
	company, manager := initSimpleDefault()
	expected := `
{  
"node": {
    "id": "1",
    "type": "COMPANY"
  },
  "labels": {},
  "children": [
    {
      "node": {
        "id": "2",
        "type": "UPLOAD_CHANNEL"
      },
      "labels": {
        "default": {}
      },
      "children": [
        {
          "node": {
            "id": "3",
            "type": "LISTING"
          },
          "labels": {
            "default": {}
          },
          "children": null
        },
        {
          "node": {
            "id": "4",
            "type": "LISTING"
          },
          "labels": {
            "default": {}
          },
          "children": null
        },
        {
          "node": {
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
}
`

	manager.CalculateBlocks(company, BlockOrder{Default, "3"})

	gojsonut.JsonCompare(t, company, expected, false)

}

func initSimpleDefault() (*Vertex, Manager) {
	company := makeCompany()

	ruler := make(Ruler)

	ruler.Add(Default,
		CurrentVertexLabelingRule{"If current node actioned default => current node has label 'default'", "", Default},
		ChildrenVertexLabelingRule{"If UPLOAD_CHANNEL actioned default => all children groups has label 'default'", "UPLOAD_CHANNEL", "", Default},
		ParentVertexLabelingRule{"If a group actioned default => it's UPLOAD_CHANNEL will have label 'default'", "", "UPLOAD_CHANNEL", Default},
		BrotherVertexLabelingRule{"If a LISTING actioned default => other listings for the same UPLOAD_CHANNEL labeled default", "", "UPLOAD_CHANNEL", "", Default},
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
