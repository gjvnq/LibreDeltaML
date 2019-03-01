package deltaml

import (
	"fmt"
	"testing"
)

func Test_TreeFromString_1(t *testing.T) {
	data := `<root param1="hi" param2="hi2">Lorem <b>Ipsum<!-- comment --><i/></b> dolor est</root>`
	tree, err := TreeFromString(data)
	if err != nil {
		fmt.Printf("error: %v", err)
		t.Error(err)
	}
	tree.Print(0)
}

func Test_TreeFromString_2(t *testing.T) {
	data := `<root param1="hi" param2="hi2">Lorem <b>Ipsum<!-- comment --></i></b> dolor est</root>`
	_, err := TreeFromString(data)
	right_err := "XML syntax error on line 1: element <b> closed by </i>"
	if err == nil || err.Error() != right_err {
		fmt.Printf("Expected error \"%s\", but got \"%+v\"", right_err, err)
		t.Fail()
	}
}
