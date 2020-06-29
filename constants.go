package gographlabel

import "errors"

var LoopInHierarchyError = errors.New("loops are not allowed in hierarchy")

const Default = "default"
