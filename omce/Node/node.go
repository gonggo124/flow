package node

type NodeType uint

const (
	THIS_WRECKED_MY_WORLD NodeType = iota // THIS IS ERROR
	PROGRAM                               // fields: Body
	INCLUDE_DECLARATION                   // fields: Body
	METHOD_DECLARATION                    // fields: Id, Body
	IDENTIFIER                            // fields: Name
	OMFN_CODE                             // fields: Code
)

type Node struct {
	Type NodeType
	Id   *Node
	Name string
	Nbt  map[string]any
	Code string
	Body []Node
}

func (n Node) String() string {
	switch n.Type {
	case PROGRAM:
		return "PROGRAM"
	case INCLUDE_DECLARATION:
		return "INCLUDE_DECLARATION"
	case METHOD_DECLARATION:
		return "METHOD DECLARATION"
	case IDENTIFIER:
		return "IDENTIFIER"
	case OMFN_CODE:
		return "OMFN CODE"
	default:
		return "THIS WRECKED MY WORLD"
	}
}

func NewNode(node_type NodeType) Node {
	new_node := Node{}
	new_node.Type = node_type
	return new_node
}
