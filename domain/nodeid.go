package domain

type NodeId struct {
	Value string
}

type NodeIdRepo interface {
	Get() (*NodeId, error)
	Put(nodeId NodeId) error
}
