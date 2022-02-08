package binarytree

import (
	"errors"
	"sync"
)

var ConnectedIndex uint64 = 0

type SingleNode interface {
	Insert(SingleNode)
	Get(interface{}) SingleNode
	Delete(interface{})
	ToggleHead()
	ToggleCanConnect()
	Length() int
	IsHead() bool
	CanConnect() bool
	All() map[interface{}]SingleNode
	Index() interface{}
}

type Tree struct {
	mutex    sync.Mutex
	nodes    map[interface{}]SingleNode
	fillNode func(interface{}, uint64) SingleNode
}

var Default = Tree{}

func (receiver *Tree) SetFillNode(function func(interface{}, uint64) SingleNode) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.fillNode = function
}

func (receiver *Tree) Insert(node interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.nodes == nil {
		receiver.nodes = make(map[interface{}]SingleNode)
	}

	receiver.nodes[node] = receiver.fillNode(node, ConnectedIndex)
	ConnectedIndex++
}

func (receiver *Tree) ToggleHead(node interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.nodes[node].ToggleHead()
}

func (receiver *Tree) ToggleCanConnect(node interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.nodes[node].ToggleCanConnect()
}

func (receiver *Tree) Get(node interface{}) SingleNode {
	return receiver.nodes[node]
}

func (receiver *Tree) Delete(node interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for nodeIndex := range receiver.nodes {
		receiver.nodes[nodeIndex].Delete(node)
	}

	delete(receiver.nodes, node)
}

func (receiver *Tree) insertConnected(parentNode interface{}, childNode SingleNode) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.nodes == nil {
		receiver.nodes = make(map[interface{}]SingleNode)
	}

	receiver.nodes[parentNode].Insert(childNode)
}

func (receiver *Tree) DeleteConnected(parentNode, childNode interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.nodes[parentNode].Delete(childNode)
}

func (receiver *Tree) All() map[interface{}]SingleNode {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.nodes
}

func (receiver *Tree) LevelNodes(level uint) []SingleNode {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	var output []SingleNode = []SingleNode{}
	for nodeIndex := range receiver.nodes {
		if receiver.nodes[nodeIndex].IsHead() {
			output = append(output, receiver.nodes[nodeIndex])
			if level == 1 {
				return output
			}
		}
	}
	if len(output) == 0 {
		return output
	}
	var index uint = 1
	var currentLevelNodes []SingleNode = output
	for {
		if len(currentLevelNodes) == 0 {
			break
		}
		output = []SingleNode{}
		for _, nodes := range currentLevelNodes {
			for indx := range nodes.All() {
				child := receiver.Get(indx)
				if child.CanConnect() {
					output = append(output, child)
				}
			}
		}
		if len(output) == 0 {
			break
		}
		if index == level-1 {
			return output
		}
		currentLevelNodes = output
		output = []SingleNode{}
		index++
	}
	return output
}

func (receiver *Tree) InsertTree(childNode interface{}) (SingleNode, error) {
	var level uint = 1
	var levelNodes []SingleNode
	for {
		levelNodes = receiver.LevelNodes(level)
		if len(levelNodes) > 0 {
			for _, node := range levelNodes {

				if node.Length() < 2 {
					receiver.insertConnected(node.Index(), receiver.Get(childNode))
					return node, nil
				}
			}
		} else {
			return nil, errors.New("no nodes to connect")
		}
		level++
	}
}

func (receiver *Tree) InsertChild(childNode interface{}, canConnect bool) (SingleNode, error) {
	if canConnect {
		receiver.ToggleCanConnect(childNode)
	}
	return receiver.InsertTree(childNode)
}
