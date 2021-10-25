package binarytree

type Tree[T any] struct {
	mutex       sync.Mutex
	nodes map[T]string
}