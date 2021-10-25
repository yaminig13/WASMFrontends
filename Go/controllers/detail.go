package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"syscall/js"
)

type product struct{
	Title string	`json:title`
	Description string	`json:description`
	Image string	`json:image`
	Price int	`json:price`
	Id string	`json:id`
	Review map[string]string	`json:review`
}

type trueVideo struct{
	video bool
	val js.Value
}

func(o *trueVideo) TOJSValue() js.Value{
	val:=js.Global().Get("Object").New()
	val.Set("video",o.video)
	return val
}
var p product
var document js.Value

func getProduct(id string) product{
	httpString:=fmt.Sprintf("http://localhost:8080/items/%v",id)
	resp,err:= http.Get(httpString)

	if err != nil {
              return p
      }

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
              err = fmt.Errorf("response status code: %d", resp.StatusCode)
              return p
      }
	var buf bytes.Buffer
      _, err = buf.ReadFrom(resp.Body)
      if err != nil {
              return p
      }
	  json.Unmarshal(buf.Bytes(), &p)
	  return p
}
// func getStream() <-chan js.Value{
	 
// 	ovideo:=trueVideo{video:true}
// 	r:=make(chan js.Value)

// 	go func(){
// 		defer close(r)
// 		r<- js.Global().Get("navigator").Get("mediaDevices").Call("getUserMedia",ovideo.TOJSValue())
// 	}()
// 	return r
// }
func startCamera(this js.Value,inputs []js.Value)interface{}{

		ovideo:=trueVideo{video:true}

	document.Call("querySelector",".cameraCanvas").Get("style").Set("display","block")
	document.Call("querySelector",".reviewImagebtn").Get("style").Set("display","none")

	video:=document.Call("getElementById","video")
	canvas:=document.Call("getElementById","canvas")
	context:=canvas.Call("getContext","2d")

	document.Call("getElementById","snap").Call("addEventListener","click",js.FuncOf(func(this js.Value,inputs []js.Value)interface{}{
		context.Call("drawImage",video,0,0,200,200)
		imageData:=canvas.Call("toDataURL")
		fmt.Println(imageData)
		return nil
	}))

	if js.Global().Get("navigator").Get("mediaDevices").Truthy()==true && js.Global().Get("navigator").Get("mediaDevices").Get("getUserMedia").Truthy()==true{
		// streamCh:=getStream()
		// stream:=<-streamCh
		// fmt.Println(stream)
		stream:=js.Global().Get("navigator").Get("mediaDevices").Call("getUserMedia",ovideo.TOJSValue())
		fmt.Println(stream.Type())
		// fmt.Println(stream.Await())
		video.Set("srcObject",stream)
		video.Call("play")
	}
	return nil
}

func displayData(prod product){
	
	document.Call("querySelector","img").Set("src",prod.Image)
	document.Call("querySelector",".title").Set("innerHTML",prod.Title)
	priceString:=fmt.Sprintf("%v€",prod.Price)
	document.Call("querySelector",".price").Set("innerHTML",priceString)
	document.Call("querySelector",".description").Set("innerHTML",prod.Description)


}

func main() {
	document=js.Global().Get("document")
	urlString := js.Global().Get("location").Get("href")
	url:=js.Global().Get("URL").New(urlString)
	id:=url.Get("searchParams").Call("get","id").String()
	prod:=getProduct(id)
	displayData(prod)
	js.Global().Get("document").Call("querySelector",".reviewImagebtn").Call("addEventListener","click",js.FuncOf(startCamera))
		select{}

}