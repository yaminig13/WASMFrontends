package service

import "syscall/js"

type Options struct{
		body string
		val js.Value
	}

	func(o *Options) TOJSValue() js.Value{
		val:=js.Global().Get("Object").New()
		val.Set("body",o.body)
		return val
	}

func service(this js.Value, inputs []js.Value) interface{} {
	
	this.Call("addEventListener","push",js.FuncOf(func(this js.Value, inputs []js.Value)interface{}{
		option:=Options{body:"Don't forget to leave your feedback! Happy Shopping!"}

		this.Get("registeration").Call("showNotification",option.TOJSValue())
		return nil
	}))
	return nil
	
}