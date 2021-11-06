package binarytree

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
)

// test structure

type MySocket struct {
	Socket           *websocket.Conn
	ID               uint64
	IsBroadcaster    bool
	Name             string
	HasStream        bool
	ConnectedSockets map[*websocket.Conn]MySocket
}

func (receiver *MySocket) Insert(node SingleNode) {
	socket := node.(*MySocket)
	receiver.ConnectedSockets[socket.Socket] = *socket
}

func (receiver *MySocket) Delete(node interface{}) {
	delete(receiver.ConnectedSockets, node.(*websocket.Conn))
}

func (receiver *MySocket) Get(node interface{}) SingleNode {
	result := receiver.ConnectedSockets[node.(*websocket.Conn)]
	return &result
}

func (receiver *MySocket) GetLength() int {
	return len(receiver.ConnectedSockets)
}

func (receiver *MySocket) IsHead() bool {
	return receiver.IsBroadcaster
}

func (receiver *MySocket) CanConnect() bool {
	return receiver.HasStream
}

func (receiver *MySocket) GetAll() map[interface{}]SingleNode {
	var output map[interface{}]SingleNode = make(map[interface{}]SingleNode)
	for indx := range receiver.ConnectedSockets {
		result := receiver.ConnectedSockets[indx]
		output[indx] = &result
	}
	return output
}

func (receiver *MySocket) ToggleHead() {
	receiver.IsBroadcaster = !receiver.IsBroadcaster
}

func (receiver *MySocket) ToggleCanConnect() {
	receiver.HasStream = !receiver.HasStream
}

func (receiver *MySocket) GetIndex() interface{} {
	return receiver.Socket
}

//\test structure

// test fill function
func fillFunction(node interface{}, socketIndex uint64) SingleNode {
	conn := node.(*websocket.Conn)
	result := MySocket{
		Socket:           conn,
		ID:               socketIndex,
		Name:             "Socket " + strconv.FormatUint(socketIndex, 10),
		ConnectedSockets: make(map[*websocket.Conn]MySocket),
	}
	return &result
}

//\test fill function

func TestWebSocketMap_Insert(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	tests := []struct {
		Input *websocket.Conn

		ExpectedLen int
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 1,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 2,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 3,
		},
	}

	for testNumber, test := range tests {
		websocketmaps.Insert(test.Input)
		if len(websocketmaps.nodes) != test.ExpectedLen {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, len(websocketmaps.nodes))
		}
	}
}

func TestWebSocketMap_Delete(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	tests := []struct {
		Input *websocket.Conn
	}{
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
	}
	for _, test := range tests {
		websocketmaps.Insert(test.Input)
	}

	for testNumber, test := range tests {
		websocketmaps.Delete(test.Input)
		deleteSocket := websocketmaps.Get(test.Input)
		if deleteSocket != nil {
			t.Errorf("Test %d : A nil was expected but got a socket", testNumber)
		}
	}
}

func TestWebSocketMap_Get(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	tests := []struct {
		Input *websocket.Conn
	}{
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
	}
	for _, test := range tests {
		websocketmaps.Insert(test.Input)
	}

	for testNumber, test := range tests {
		mySocket := websocketmaps.Get(test.Input)
		if mySocket.(*MySocket).Socket != test.Input {
			t.Errorf("Test %d : A socket was expected but got another one", testNumber)
		}
	}
}

func TestWebSocketMap_InsertConnected(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	parent := &websocket.Conn{}
	websocketmaps.Insert(parent)
	tests := []struct {
		Input *websocket.Conn

		ExpectedLen int
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 1,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 2,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 3,
		},
	}

	for testNumber, test := range tests {
		websocketmaps.Insert(test.Input)
		websocketmaps.InsertConnected(parent, websocketmaps.Get(test.Input))
		parentSocket := websocketmaps.Get(parent)
		if parentSocket.GetLength() != test.ExpectedLen {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, parentSocket.GetLength())
		}
	}
}

func TestWebSocketMap_DeleteConnected(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)
	parent := &websocket.Conn{}
	websocketmaps.Insert(parent)
	tests := []struct {
		Input *websocket.Conn

		ExpectedLen int
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 0,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 0,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 0,
		},
	}

	for testNumber, test := range tests {
		websocketmaps.Insert(test.Input)
		websocketmaps.InsertConnected(parent, websocketmaps.Get(test.Input))
		parentSocket := websocketmaps.Get(parent)
		if parentSocket.GetLength() != test.ExpectedLen+1 {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen+1, parentSocket.GetLength())
		}
		websocketmaps.DeleteConnected(parent, test.Input)
		parentSocket = websocketmaps.Get(parent)
		if parentSocket.GetLength() != test.ExpectedLen {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, parentSocket.GetLength())
		}
	}
}

func TestWebSocketMap_GetAll(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	tests := []struct {
		Input *websocket.Conn

		ExpectedLen int
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 1,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 2,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 3,
		},
	}

	for testNumber, test := range tests {
		websocketmaps.Insert(test.Input)
		allNodes := websocketmaps.GetAll()
		if len(allNodes) != test.ExpectedLen {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, len(allNodes))
		}
	}
}

func TestWebSocketMap_LevelNodes(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	var broadcaster *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(broadcaster)
	websocketmaps.ToggleHead(broadcaster)
	websocketmaps.ToggleCanConnect(broadcaster)
	levelNodes := websocketmaps.LevelNodes(1)
	if len(levelNodes) != 1 {
		t.Errorf("Test with a broadcaster should return 1 nodes in level 1 but it retuens %d", len(levelNodes))
	}
	var nodeOne *websocket.Conn = &websocket.Conn{}
	var nodeTwo *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeOne)
	websocketmaps.ToggleCanConnect(nodeOne)
	websocketmaps.InsertConnected(broadcaster, websocketmaps.Get(nodeOne))
	websocketmaps.Insert(nodeTwo)
	websocketmaps.ToggleCanConnect(nodeTwo)
	websocketmaps.InsertConnected(broadcaster, websocketmaps.Get(nodeTwo))
	levelNodes = websocketmaps.LevelNodes(2)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the broadcaster should return 2 nodes in level 2 but it retuens %d", len(levelNodes))
	}
	var nodeThree *websocket.Conn = &websocket.Conn{}
	var nodeFour *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeThree)
	websocketmaps.ToggleCanConnect(nodeThree)
	websocketmaps.InsertConnected(nodeOne, websocketmaps.Get(nodeThree))
	websocketmaps.Insert(nodeFour)
	websocketmaps.ToggleCanConnect(nodeFour)
	websocketmaps.InsertConnected(nodeOne, websocketmaps.Get(nodeFour))
	levelNodes = websocketmaps.LevelNodes(3)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the nodeOne should return 2 nodes in level 3 but it retuens %d", len(levelNodes))
	}
}

func TestWebSocketMap_InsertTree(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	var broadcaster *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(broadcaster)
	websocketmaps.ToggleHead(broadcaster)
	websocketmaps.ToggleCanConnect(broadcaster)
	levelNodes := websocketmaps.LevelNodes(1)
	if len(levelNodes) != 1 {
		t.Errorf("Test with a broadcaster should return 1 nodes in level 1 but it retuens %d", len(levelNodes))
	}
	var nodeOne *websocket.Conn = &websocket.Conn{}
	var nodeTwo *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeOne)
	websocketmaps.ToggleCanConnect(nodeOne)
	_, err := websocketmaps.InsertTree(nodeOne)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	websocketmaps.Insert(nodeTwo)
	websocketmaps.ToggleCanConnect(nodeTwo)
	_, err = websocketmaps.InsertTree(nodeTwo)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	levelNodes = websocketmaps.LevelNodes(2)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the broadcaster should return 2 nodes in level 2 but it retuens %d", len(levelNodes))
	}
	var nodeThree *websocket.Conn = &websocket.Conn{}
	var nodeFour *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeThree)
	websocketmaps.ToggleCanConnect(nodeThree)
	_, err = websocketmaps.InsertTree(nodeThree)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	websocketmaps.Insert(nodeFour)
	websocketmaps.ToggleCanConnect(nodeFour)
	_, err = websocketmaps.InsertTree(nodeFour)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	levelNodes = websocketmaps.LevelNodes(3)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the nodeOne should return 2 nodes in level 3 but it retuens %d", len(levelNodes))
	}
}

func TestWebSocketMap_InsertChild(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	var broadcaster *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(broadcaster)
	websocketmaps.ToggleHead(broadcaster)
	websocketmaps.ToggleCanConnect(broadcaster)
	levelNodes := websocketmaps.LevelNodes(1)
	if len(levelNodes) != 1 {
		t.Errorf("Test with a broadcaster should return 1 nodes in level 1 but it retuens %d", len(levelNodes))
	}
	var level uint = 1
	fmt.Println('-')
	for {
		fmt.Println("*Level ", level)
		levelNodes = websocketmaps.LevelNodes(level)
		if len(levelNodes) == 0 {
			break
		}
		for i, testNode := range levelNodes {
			fmt.Println("*test node ", i, testNode.(*MySocket).ID)
		}
		level++
	}
	var nodeOne *websocket.Conn = &websocket.Conn{}
	var nodeTwo *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeOne)
	_, err := websocketmaps.InsertChild(nodeOne, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	websocketmaps.Insert(nodeTwo)
	_, err = websocketmaps.InsertChild(nodeTwo, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	levelNodes = websocketmaps.LevelNodes(2)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the broadcaster should return 2 nodes in level 2 but it retuens %d", len(levelNodes))
	}
	fmt.Println('-')
	level = 1
	for {
		fmt.Println("*Level ", level)
		levelNodes = websocketmaps.LevelNodes(level)
		if len(levelNodes) == 0 {
			break
		}
		for i, testNode := range levelNodes {
			fmt.Println("*test node ", i, testNode.(*MySocket).ID)
		}
		level++
	}
	var nodeThree *websocket.Conn = &websocket.Conn{}
	var nodeFour *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeThree)
	_, err = websocketmaps.InsertChild(nodeThree, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	websocketmaps.Insert(nodeFour)
	_, err = websocketmaps.InsertChild(nodeFour, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	levelNodes = websocketmaps.LevelNodes(3)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the nodeOne should return 2 nodes in level 3 but it retuens %d", len(levelNodes))
	}
	fmt.Println('-')
	level = 1
	for {
		fmt.Println("*Level ", level)
		levelNodes = websocketmaps.LevelNodes(level)
		if len(levelNodes) == 0 {
			break
		}
		for i, testNode := range levelNodes {
			fmt.Println("*test node ", i, testNode.(*MySocket).ID)
		}
		level++
	}
}

func TestWebSocketMap_3Aud(t *testing.T) {
	var websocketmaps Tree = Tree{}
	websocketmaps.SetFillNode(fillFunction)

	var broadcaster *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(broadcaster)
	websocketmaps.ToggleHead(broadcaster)
	websocketmaps.ToggleCanConnect(broadcaster)
	levelNodes := websocketmaps.LevelNodes(1)
	if len(levelNodes) != 1 {
		t.Errorf("Test with a broadcaster should return 1 nodes in level 1 but it retuens %d", len(levelNodes))
	}
	var nodeOne *websocket.Conn = &websocket.Conn{}
	var nodeTwo *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeOne)
	_, err := websocketmaps.InsertChild(nodeOne, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	websocketmaps.Insert(nodeTwo)
	_, err = websocketmaps.InsertChild(nodeTwo, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	levelNodes = websocketmaps.LevelNodes(2)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the broadcaster should return 2 nodes in level 2 but it retuens %d", len(levelNodes))
	}

	var nodeThree *websocket.Conn = &websocket.Conn{}
	var nodeFour *websocket.Conn = &websocket.Conn{}
	websocketmaps.Insert(nodeThree)
	_, err = websocketmaps.InsertChild(nodeThree, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	websocketmaps.Insert(nodeFour)
	_, err = websocketmaps.InsertChild(nodeFour, true)
	if err != nil {
		t.Errorf("Error Happend! %e", err)
	}
	levelNodes = websocketmaps.LevelNodes(3)
	if len(levelNodes) != 2 {
		t.Errorf("Test with 2 nodes connected to the nodeOne should return 2 nodes in level 3 but it retuens %d", len(levelNodes))
	}
}
