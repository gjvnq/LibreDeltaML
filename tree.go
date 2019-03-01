package deltaml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Tree struct {
	Token    xml.Token
	Children []*Tree
	EndToken xml.Token
}

func (tree *Tree) AddChild(child *Tree) {
	if tree.Children == nil {
		tree.Children = make([]*Tree, 0)
	}
	tree.Children = append(tree.Children, child)
}

func (tree Tree) ToXML() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := xml.NewEncoder(buf)
	err := tree.ToXMLWithEncoder(encoder)
	if err != nil {
		return nil, err
	}
	encoder.Flush()
	return buf.Bytes(), nil
}

func (tree Tree) ToXMLWithEncoder(encoder *xml.Encoder) error {
	err := encoder.EncodeToken(tree.Token)
	if err != nil {
		panic(err)
		return err
	}
	for _, child := range tree.Children {
		err = child.ToXMLWithEncoder(encoder)
		if err != nil {
			panic(err)
			return err
		}
	}
	if tree.EndToken != nil {
		err = encoder.EncodeToken(tree.EndToken)
		if err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

func (tree Tree) Print(level int) {
	tabs := ""
	for i := 0; i < level; i++ {
		tabs += "\t"
	}

	switch tk := tree.Token.(type) {
	case xml.CharData:
		fmt.Printf(tabs+"Token: [CharData] \"%s\"\n", string(tk))
	case xml.Comment:
		fmt.Printf(tabs+"Token: [Comment] \"%s\"\n", string(tk))
	default:
		fmt.Printf(tabs+"Token: %+v\n", tree.Token)
	}
	fmt.Printf(tabs+"Children: %d\n", len(tree.Children))
	for _, child := range tree.Children {
		child.Print(level + 1)
	}
	fmt.Printf(tabs+"EndToken: %+v\n", tree.EndToken)
}

type TreeStack struct {
	Data []*Tree
	Top  int
}

func (stack *TreeStack) Init() {
	stack.Data = make([]*Tree, 0)
}

func (stack *TreeStack) Push(tree *Tree) {
	stack.Data = append(stack.Data, tree)
	stack.Top++
}

func (stack *TreeStack) Peek() *Tree {
	if stack.Top == 0 {
		return nil
	}
	return stack.Data[stack.Top-1]
}

func (stack *TreeStack) Pop() *Tree {
	ans := stack.Peek()
	stack.Top--
	return ans
}

func TreeFromString(data string) (Tree, error) {
	decoder := xml.NewDecoder(strings.NewReader(data))

	// Prepare variables
	var err error
	var token xml.Token
	tree_stack := TreeStack{}
	tree_stack.Init()

	// Create root
	root := Tree{}
	token, err = decoder.Token()
	switch tk := token.(type) {
	case xml.StartElement:
		root.Token = tk
		tree_stack.Push(&root)
	default:
		return Tree{}, errors.New("First XML token MUST be a start element")
	}

	// For each token
	for err == nil {
		// Get the current sub tree
		cursor := tree_stack.Peek()
		// Get the next token
		token, err = decoder.Token()
		// Copy the token we got into a new sub tree
		child := Tree{}
		child.Token = xml.CopyToken(token)

		// Decide what to do with the token
		switch token.(type) {
		case xml.StartElement:
			// Add sub tree and go down one level
			cursor.AddChild(&child)
			tree_stack.Push(&child)
		case xml.EndElement:
			// Go up one level and annotate if that sub tree has a dedicated end tag
			tree_stack.Pop().EndToken = xml.CopyToken(token)
		default:
			// If its not nil, just copy it
			if token != nil {
				cursor.AddChild(&child)
			}
		}
	}
	if err != nil && err != io.EOF {
		return root, err
	}

	return root, nil
}
