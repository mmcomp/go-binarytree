package binarytree

import (
	"fmt"
	"sync"
)

var ConnectedIndex uint64 = 0

type SingleNode interface {
	Insert(interface{}, uint64)
	Get(interface{}) SingleNode
	Delete(interface{})
	ToggleHead()
	ToggleCanConnect()
	GetLength() int
	IsHead() bool
	CanConnect() bool
	GetAll() map[interface{}]SingleNode
}

type Tree struct {
	mutex    sync.Mutex
	nodes    map[interface{}]SingleNode
	fillNode func(interface{}, uint64) SingleNode
}

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

	fmt.Println("Toggling a head, corrent head is ", receiver.nodes[node].IsHead())
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

func (receiver *Tree) InsertConnected(parentNode, childNode interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.nodes == nil {
		receiver.nodes = make(map[interface{}]SingleNode)
	}

	receiver.nodes[parentNode].Insert(childNode, ConnectedIndex)
	ConnectedIndex++
}

func (receiver *Tree) DeleteConnected(parentNode, childNode interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.nodes[parentNode].Delete(childNode)
}

func (receiver *Tree) GetAll() map[interface{}]SingleNode {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.nodes
}

func (receiver *Tree) LevelNodes(level uint) []SingleNode {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	fmt.Println("LevelNodes ", level)
	var output []SingleNode = []SingleNode{}
	for nodeIndex := range receiver.nodes {
		if receiver.nodes[nodeIndex].IsHead() {
			output = append(output, receiver.nodes[nodeIndex])
			fmt.Println("Head found !")
			if level == 1 {
				fmt.Println("Requested level is 1 so end")
				return output
			}
		}
	}
	if len(output) == 0 {
		fmt.Println("Head did not found!")
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
			for _, child := range nodes.GetAll() {
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
