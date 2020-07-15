package gographlabel

func DefaultRuler() Ruler {
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
	return ruler
}
