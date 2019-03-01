package deltaml

import (
	"fmt"
	"testing"
)

func Test_TreeFromString_1(t *testing.T) {
	data := `<root param1="hi" param2="hi2">Lorem <b>Ipsum</b> dolor <!-- comment --> est<i/></root>`
	tree, err := TreeFromString(data)
	if err != nil {
		fmt.Printf("error: %v", err)
		t.Error(err)
	}
	fmt.Printf("tree: %+v\n", tree)
}
