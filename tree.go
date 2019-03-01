package deltaml

import (
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
		// fmt.Printf("Begin element %s:%s\n", tk.Name.Space, tk.Name.Local)
		root.Token = tk
		tree_stack.Push(&root)
	default:
		return Tree{}, errors.New("First XML token MUST be a start element")
	}

	// For each token
	for err == nil {
		fmt.Printf("     ---\n")
		cursor := tree_stack.Peek()
		token, err = decoder.Token()
		child := Tree{}
		child.Token = xml.CopyToken(token)
		child.Print(0)

		switch tk := token.(type) {
		case xml.StartElement:
			fmt.Printf(">>>> Begin element %s:%s\n", tk.Name.Space, tk.Name.Local)
			cursor.AddChild(&child)
			tree_stack.Push(&child)
		case xml.CharData:
			fmt.Printf(">>>> Char data: %s\n", string(tk))
			cursor.AddChild(&child)
		case xml.EndElement:
			fmt.Printf(">>>> End element %s:%s\n", tk.Name.Space, tk.Name.Local)
			tree_stack.Pop().EndToken = xml.CopyToken(token)
		default:
			fmt.Printf(">>>> Got comment, directive or something: %+v\n", tk)
			if token != nil {
				cursor.AddChild(&child)
			}
		}
		root.Print(0)
	}
	if err != nil && err != io.EOF {
		fmt.Printf(">>>> Got error: %+v\n", err)
		return root, err
	}
	fmt.Printf("-----END-----\n")
	root.Print(0)
	fmt.Printf("---------\n")

	return root, nil
}
