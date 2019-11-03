package main

import (
	"bytes"
	"fmt"
)

type Cat struct {
	Name  string
	Type  string
	Color string
}

func initLedger() {
	cats := []Cat{
		Cat{Name: "Da Hua", Type: "American Shorthair", Color: "White"},
		Cat{Name: "Xiao Huang", Type: "Ragdoll", Color: "Yellow"},
		Cat{Name: "Xiao Mi", Type: "Leopard Cat", Color: "Black"},
		Cat{Name: "Lai Fu", Type: "Garfield", Color: "Gray"},
		Cat{Name: "Mi Fan", Type: "Persian Cat", Color: "White,Yellow"},
	}

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	var i int
	for i = 0; i < len(cats); i++ {
		buffer.WriteString("{name: " + cats[i].Name + "},\n")
		// fmt.Printf("%s \n",cats[i])
	}
	buffer.WriteString("]")
	fmt.Printf("queryAllCats:\n%s\n", buffer.String())
}

func main() {
	initLedger()
}
