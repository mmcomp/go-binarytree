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

func (receiver *MySocket) Insert(node interface{}, socketIndex uint64) {
	receiver.ConnectedSockets[node.(*websocket.Conn)] = MySocket{
		Socket:           node.(*websocket.Conn),
		ID:               socketIndex,
		IsBroadcaster:    false,
		ConnectedSockets: make(map[*websocket.Conn]MySocket),
	}
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
		websocketmaps.InsertConnected(parent, test.Input)
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
		websocketmaps.InsertConnected(parent, test.Input)
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

	tests := []struct {
		Input       *websocket.Conn
		ExpectedLen int
		IsHead      bool
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 1,
			IsHead:      true,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 2,
			IsHead:      false,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 3,
			IsHead:      false,
		},
	}

	for _, test := range tests {
		websocketmaps.Insert(test.Input)
		if test.IsHead {
			websocketmaps.ToggleHead(test.Input)
		}
	}
	// for testNumber, test := range tests {
	levelNodes := websocketmaps.LevelNodes(1)
	fmt.Println("len ", len(levelNodes))
	// allNodes := websocketmaps.GetAll()
	// fmt.Println("all nodes ", allNodes)
	// for indx := range allNodes {
	// 	fmt.Println("Is Head ? ", allNodes[indx].IsHead())
	// }
	// if len(allNodes) != test.ExpectedLen {
	// 	t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, len(allNodes))
	// }
	// }
}
