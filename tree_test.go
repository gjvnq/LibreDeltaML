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

func Test_ToXML(t *testing.T) {
	data := `<a:root param1="hi" param2="hi2">Lorem <b:b>Ipsum<!-- comment --><c:i/></b:b> dolor est</a:root>`
	right_ans := `<root xmlns="a" param1="hi" param2="hi2">Lorem <b xmlns="b">Ipsum<!-- comment --><i xmlns="c"></i></b> dolor est</root>`
	tree, err := TreeFromString(data)
	if err != nil {
		t.Error(err)
	}
	ans, err := tree.ToXML()
	if err != nil {
		t.Error(err)
	}
	if string(ans) != right_ans {
		t.Log("Got         : " + string(ans))
		t.Log("But expected: " + right_ans)
		t.Fail()
	}
}
