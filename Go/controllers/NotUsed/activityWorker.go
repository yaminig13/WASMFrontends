package main

import (
	"fmt"
	"syscall/js"
)


var addClicks int
var removeClicks int

func main() {
	onmessage := func(this js.Value, inputs []js.Value) interface{} {
		if inputs[0].Get("data").String() == "add" {
			addClicks += 1
			fmt.Println(removeClicks)

		}
		if inputs[0].Get("data").String() == "remove" {
			removeClicks += 1
			fmt.Println(removeClicks)
		}
		return nil
	}
	_=onmessage

	// js.Global().Call("setInterval",js.FuncOf(func(){
	// 	msg:=""
	// 	js.Global().Call("postMessage","")
	// }))
}