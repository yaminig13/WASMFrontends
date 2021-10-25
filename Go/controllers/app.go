package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"syscall/js"
	"time"
)

var searchHistory string
var cartItems []products
var activityWorker js.Value
var cartTotal int

type activityLog struct{
	Type string 	`json:type`
	TimeStamp string 	`json:timeStamp`	
}
var logItems []activityLog

type DataStore struct{}
var product DataStore
var store DataStore

type products struct{
	Title string	`json:title`
	Description string	`json:description`
	Image string	`json:image`
	Price int	`json:price`
	Id string	`json:id`
	Review map[string]string	`json:review`
}
var p [5]products

type stores struct{
	Latitude float64	`json:latitude`
	Longitude float64	`json:longitude`
	Name string	`json:Name`
	Address string	`json:Address`
}
var s [3]stores

// fetch stores
func (store *DataStore) getStores() []stores{
	resp,err:= http.Get("http://localhost:8080/stores/")

	if err != nil {
              return nil
      }

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
              err = fmt.Errorf("response status code: %d", resp.StatusCode)
              return nil
      }
	var buf bytes.Buffer
      _, err = buf.ReadFrom(resp.Body)
      if err != nil {
              return nil
      }
	  json.Unmarshal(buf.Bytes(), &s)
	  return s[:]
}

// fetch products
func (product *DataStore) getProducts()  []products{
	// trial by fetch method

	// response:=js.Global().Call("fetch","http://localhost:8080/items");
	// fmt.Println(response.String())

	resp,err:= http.Get("http://localhost:8080/items/")

	if err != nil {
              return nil
      }

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
              err = fmt.Errorf("response status code: %d", resp.StatusCode)
              return nil
      }
	var buf bytes.Buffer
      _, err = buf.ReadFrom(resp.Body)
      if err != nil {
              return nil
      }
	  json.Unmarshal(buf.Bytes(), &p)

	  
	// trial by xhttp method

		// xhttp:=js.Global().Get("XMLHttpRequest").New()
		// xhttp.Call("open","GET","http://localhost:8080/items",true)

		// xhttp.Set("onload",js.FuncOf(
		// func(this js.Value,inputs []js.Value) interface{}{
			
				
		// 	if this.Get("status").Int()<300 && this.Get("status").Int()>=200{
		// 		response:=this.Get("responseText").String()
		// 		var rgx=regexp.MustCompile(`(\{.*?\])`)
		// 		rs:=rgx.FindAllStringSubmatch(response,-1)

		// 		for i:=0;i<len(p);i++ {
		// 			rs[i][0]=rs[i][0]+"}"
		// 			bytes:=[]byte(rs[i][0])
		// 			json.Unmarshal(bytes,&p[i])
		// 		}
		// 	}
		// 	return nil
		// }))
		// xhttp.Call("send")
return p[:]
} 



type UI struct{}
var ui UI

// web component cart-element
func cartElementSetup(cartElement js.Value,cartItem products){
	template:=`<img  alt="Cartimage">
    <div>
        <h4></h4>
        <h5>€</h5>
        <span class="remove-item">Remove</span>
    </div>`

	cartElement.Set("innerHTML",template)
	priceString:=fmt.Sprintf("%v€",cartItem.Price)
	cartElement.Call("querySelector","img").Call("setAttribute","src",cartItem.Image)
	cartElement.Call("querySelector","h4").Set("innerHTML",cartItem.Title)
	cartElement.Call("querySelector","h5").Set("innerHTML",priceString)
	cartElement.Call("querySelector",".remove-item").Call("setAttribute","data-id",cartItem.Id)

	// deleting item from cart
	cartElement.Call("querySelector",".remove-item").Call("addEventListener","click",js.FuncOf(
		func(this js.Value,inputs []js.Value) interface{}{

			activityWorker.Call("postMessage","remove")
			
			cartElement.Call("remove")
			cartTotal=cartTotal-cartItem.Price
			js.Global().Get("document").Call("querySelector",".cart-total").Set("innerHTML",cartTotal)	

			for i,c:=range cartItems{
				if c.Id==cartItem.Id{
					cartItems=append(cartItems[:i],cartItems[i+1:]...)
					break
				}
			}
			p,_:=json.Marshal(cartItems)
			js.Global().Get("localStorage").Call("setItem","cart",string(p))

			return nil
		}))
}

func (ui *UI) populateCart(){

	// fetch cart elements from local storage
	cartTotal=0
	js.Global().Get("document").Call("querySelector",".cart-content").Set("innerHTML","")
	js.Global().Get("document").Call("querySelector",".cart-total").Set("innerHTML",cartTotal)
	localCart:=js.Global().Get("localStorage").Call("getItem","cart").String()
	cartItemsByte:=[]byte(localCart)
	json.Unmarshal(cartItemsByte, &cartItems)


	for _,cartItem:=range cartItems{
		div:=js.Global().Get("document").Call("createElement","div")
		div.Get("classList").Call("add","cart-item")
		cartElement:=js.Global().Get("document").Call("createElement","cart-element")

		// call cart web component

		cartElementSetup(cartElement,cartItem)
		div.Call("appendChild",cartElement)
		cartTotal+=cartItem.Price
		js.Global().Get("document").Call("querySelector",".cart-content").Call("appendChild",div)	

	}
	js.Global().Get("document").Call("querySelector",".cart-total").Set("innerHTML",cartTotal)	

}

func (ui *UI) setupCart(){
	cartOverlay:=js.Global().Get("document").Call("querySelector",".cart-overlay")
	cartDOM:=js.Global().Get("document").Call("querySelector",".cart")
	cartBtn:=js.Global().Get("document").Call("querySelector",".cart-btn")
	closeCartBtn:=js.Global().Get("document").Call("querySelector",".close-cart")

	// show Cart
	cartBtn.Call("addEventListener","click",js.FuncOf(func(this js.Value,inputs []js.Value)interface{}{
		ui.populateCart();
		cartOverlay.Get("classList").Call("add","transparentBcg")
		cartDOM.Get("style").Set("display","inline")
		return nil
	}))

	// hide Cart
	closeCartBtn.Call("addEventListener","click",js.FuncOf(
		func(this js.Value,inputs []js.Value)interface{}{
			cartOverlay.Get("classList").Call("remove","transparentBcg")
			cartDOM.Get("style").Set("display","none")
			return nil
	}))
}


func (ui *UI) searchProducts(searchTerm string){
	// search products from indexdb
	searchTerm=strings.ToLower(searchTerm)

	dbName:="GoLang"

	request:=js.Global().Get("indexedDB").Call("open",dbName,1);

	request.Set("onerror",js.FuncOf(
		func (this js.Value,inputs []js.Value) interface{} {
			fmt.Println("error found")
			return nil
		}))
	
	request.Set("onsuccess",js.FuncOf(
		func (this js.Value,inputs []js.Value) interface{} {
			var result []products
			db:=inputs[0].Get("target").Get("result");
			txn:=db.Call("transaction","products","readwrite")
			tempProducts:=txn.Call("objectStore","products")
			b:=[]byte(tempProducts.String())
	  		json.Unmarshal(b, &result)
			fmt.Println(result)
			titleIndex:=tempProducts.Call("index","Title")
			query:=titleIndex.Call("getAll",searchTerm)

			query.Set("onsuccess",js.FuncOf(
				func (this js.Value,inputs []js.Value) interface{}  {
					b:=[]byte(query.Get("result").String())
	  				json.Unmarshal(b, &result)
					if len(result)==0{
						fmt.Println("no product found")
					} else{
						// ui.displayProducts(result)
					}
					return nil
				}))
			return nil
		}))
}

func customEl(this js.Value, i []js.Value) interface{} {
		return nil
	}

func getCartBtns(this js.Value, respProduct products){

	// add to cart by adding in local storage
	flag:=false
	button:=this.Call("querySelector",".bag-btn")
	button.Call("addEventListener","click",js.FuncOf(func(this js.Value, i []js.Value) interface{}{
		activityWorker.Call("postMessage","add")
		localCart:=js.Global().Get("localStorage").Call("getItem","cart").String()
		cartItemsByte:=[]byte(localCart)
	  	json.Unmarshal(cartItemsByte, &cartItems)
		
			for _,c:=range cartItems{
				if c.Id==respProduct.Id{
					flag=true
				}
			}
			if flag==true{
				js.Global().Get("alert").Invoke("already in cart!")
			}else{
				cartItems=append(cartItems,respProduct)
				p,_:=json.Marshal(cartItems)
				js.Global().Get("localStorage").Call("setItem","cart",string(p))
				js.Global().Get("alert").Invoke("Successfully added to cart!")
			}
		
		return nil
	}))
}

func (ui *UI) displayViewedProducts(){
	ca:=strings.Split(js.Global().Get("document").Get("cookie").String(),";")
	if len(ca)>0{
		js.Global().Get("document").Call("querySelector",".viewedProducts").Get("style").Set("display","block")
	}

	for i := 0; i < len(ca); i++ {
		val:=(strings.Split(ca[i],"="))[1]
		imgDiv:=js.Global().Get("document").Call("createElement","img")

		imgDiv.Call("setAttribute","src",val)
		js.Global().Get("document").Call("querySelector",".viewed-product").Call("appendChild",imgDiv)
	}

}

func navigate(this js.Value,respProduct products) interface{}{
	// fmt.Println(this.Call("querySelector","h3"))
	this.Call("querySelector","img").Call("addEventListener","click",js.FuncOf(func(this js.Value, i []js.Value) interface{}{
		
		urlString:=fmt.Sprintf("http://localhost:8080/detail.html?id=%v",respProduct.Id)
		js.Global().Set("location",urlString)

		// Cookies
		
		cookie:=js.Global().Get("document").Get("cookie")
		cookieString:=cookie.String()
		
		ca:=strings.Split(cookieString,";")

		cookieString=fmt.Sprintf("%v=%v",respProduct.Title,respProduct.Image)
		js.Global().Get("document").Set("cookie",cookieString)

		if len(ca)>2{
			cookieString=fmt.Sprintf("%v=;expires=Thu, 01 Jan 1970 00:00:00 GMT",ca[0])
			js.Global().Get("document").Set("cookie",cookieString)
		}

		fmt.Println(cookie)
		ui.displayViewedProducts()
		
		return nil
	}))

	return nil
}


func (ui *UI) displayProducts(respProducts []products){
// display fetched products via web components
	for _, respProduct := range respProducts {
		respProduct.Description=""
		productelement:=js.Global().Get("document").Call("createElement","product-element")
		template:=`<article class="product">
    	<div class="img-container">
        <img height=10 width=10 class="product-img">
        <button class="bag-btn" >
            <i class="fas fa-shopping-cart"></i>
            Add to cart
        </button>
    	</div>
    	<h3></h3>
    	<h4></h4>
		</article>`
		productelement.Set("innerHTML",template)
		productelement.Call("querySelector",".product-img").Set("src",respProduct.Image)
		productelement.Call("querySelector","button").Call("setAttribute","data-id",respProduct.Id)	
		productelement.Call("querySelector","h3").Set("innerHTML",respProduct.Title)
		priceString:=fmt.Sprintf("%v€",respProduct.Price)
		productelement.Call("querySelector","h4").Set("innerHTML",priceString)

		productelement.Call("querySelector",".product").Call("addEventListener","mouseover",js.FuncOf(
			func(this js.Value, i []js.Value) interface{}{
				productelement.Call("querySelector",".bag-btn").Get("style").Set("display","flex")
				return nil
			}))
		
		productelement.Call("querySelector",".product").Call("addEventListener","mouseleave",js.FuncOf(
			func(this js.Value, i []js.Value) interface{}{
				productelement.Call("querySelector",".bag-btn").Get("style").Set("display","none")
				return nil
			}))	

		getCartBtns(productelement,respProduct)
		// navigation
		navigate(productelement,respProduct)
		// productelement.Call("querySelector","img").Call("addEventListener","click",js.FuncOf(navigate))
  		js.Global().Get("document").Call("querySelector",".products-center").Call("appendChild",productelement)	

	}
}

func (ui *UI) filterFunction(this js.Value,inputs []js.Value)interface{} {

	// filter the search terms and retrieve from session storage
	if (js.Global().Get("sessionStorage").Call("getItem","searchHistory").String()!="") {
		var result string
		var sessionTerms []string

		filterTemp:=inputs[0].Get("target").Get("value")
		filter:=strings.ToLower(filterTemp.String())
		sessionTermsString:=js.Global().Get("sessionStorage").Call("getItem","searchHistory").String()
		sessionTerms=strings.Split(sessionTermsString,",")
		for _, term:=range sessionTerms{
			if strings.Contains(term,filter){
				result+=fmt.Sprintf("<option value=%v>",term)
			}
		}
		js.Global().Get("document").Call("querySelector","#productSuggestion").Set("innerHTML",result)

	}
	return nil
}


type Storage struct{}
var storage Storage


func (storage *Storage) setSessionStorage(searchTerm string){

	// add searched terms into session storage 
	sessionStorage:= js.Global().Get("sessionStorage").Call("getItem","searchHistory").String()
	flag:=false
	if js.Global().Get("Storage").Truthy()==true{
		if sessionStorage!=""{
			searchHistory=sessionStorage
			sessionTerms:=strings.Split(sessionStorage,",")

			for _,term:=range sessionTerms{
				if term==searchTerm{
					flag=true
				}
			}
			if flag==false{
				searchHistory+=fmt.Sprintf(",%v",searchTerm)
			}
		} else{
			searchHistory=searchTerm
		}
		js.Global().Get("sessionStorage").Call("setItem","searchHistory",searchHistory)
	}
}


// variables for index db setup
// type keyPairs struct{
// 		keyPath string
// 		val js.Value
// 	}

// func(o *keyPairs) TOJSValue() js.Value{
// 	val:=js.Global().Get("Object").New()
// 	val.Set("keyPath",o.keyPath)
// 	return val
// }
type falseIndex struct{
	unique bool
	val js.Value
}

func(o *falseIndex) TOJSValue() js.Value{
	val:=js.Global().Get("Object").New()
	val.Set("unique",o.unique)
	return val
}

func (storage *Storage) createIndexDB(respProducts []products)  {
	
	// setting up of index db for searching of products
	// okeyPath:=keyPairs{keyPath:"Id"}
	ounique:=falseIndex{unique:true}
	dbName:="GoLang"
	request:=js.Global().Get("indexedDB").Call("open",dbName,1)
	request.Set("onerror",js.FuncOf(
		func (this js.Value,inputs []js.Value) interface{} {
			fmt.Println("error found")
			return nil
		}))
	request.Set("onupgradeneeded",js.FuncOf(
		func (this js.Value,inputs []js.Value) interface{} {
			db:=inputs[0].Get("target").Get("result");
			objectStore:=db.Call("createObjectStore","products")
			objectStore.Call("createIndex","Title","Title",ounique.TOJSValue())
			objectStore.Get("transaction").Set("oncomplete",js.FuncOf(
				func(this js.Value,inputs []js.Value) interface{}{
					productStore:=db.Call("transaction","products","readwrite").Call("objectStore","products")
					for _,respProduct:=range respProducts{
						p,_:=json.Marshal(respProduct)
						productStore.Call("add",string(p),respProduct.Id)
					}
					return nil
				}))
		return nil
		}))
}


type API struct{}
var api API
var respStores []stores
type Options struct{
		body string
		val js.Value
	}

	func(o *Options) TOJSValue() js.Value{
		val:=js.Global().Get("Object").New()
		val.Set("body",o.body)
		return val
	}




func (api *API) setupPushAPI(){
	fmt.Println(js.Global().Get("PushManager"))
	sw:=js.Global().Get("navigator").Get("serviceWorker").Call("register","controllers/serviceWorker.js")
	fmt.Println(sw)

	fmt.Println(sw.Get("pushManager"))
	register:=func() <-chan js.Value{
		r:=make(chan js.Value)
		
		go func(){
			fmt.Println("in")
			fmt.Println(sw)
			defer close(r)
			time.Sleep(time.Second * 3)
			r<-sw.Get("ready")
		}()
		// s:=r.Get("pushManager").Call("getSubscription")		
		return r
	}
	subsCh:=register()
	subs:=<-subsCh
	fmt.Println(subs)
}


func (api *API) notifiyUsers(){
	
	option:=Options{body:"Sale on all products for limited time!"}

	if js.Global().Get("Notification").Truthy()==true &&  js.Global().Get("Notification").Get("perrmission").String()!="denied"{
		js.Global().Get("Notification").Call("requestPermission",js.FuncOf(func (this js.Value,inputs []js.Value) interface{}  {
			js.Global().Get("Notification").New("Good News",option.TOJSValue())
			return nil
		}))
	}
}

func findStore(this js.Value,inputs []js.Value) interface{}{
	// search store based on user location
	userLat:=inputs[0].Get("coords").Get("latitude").Float()
	userLong:=inputs[0].Get("coords").Get("longitude").Float()
	
	minDif:=999999
	var closest int
	
	for index := 0; index < len(respStores); index++ {
		var lat1=userLat*math.Pi/180;
        var lon1=userLong*(math.Pi/180);
        var lat2=respStores[index].Latitude*math.Pi/180;
        var lon2=respStores[index].Longitude*math.Pi/180;

        var R=6371.0;
        var x=(lon2-lon1)*math.Cos((lat1+lat2)/2);
        var y=(lat2-lat1);
        var dif=math.Sqrt(x*x+y*y)*R;

		if dif<float64(minDif){
			closest=index
			minDif=int(dif)
		}
	}
	js.Global().Get("document").Call("querySelector",".location").Set("innerHTML",respStores[closest].Address)
	
	return nil
}

func (api *API) getLocation(this js.Value,inputs []js.Value) interface{}  {
	// get user location
	if js.Global().Get("navigator").Get("geolocation").Truthy()==true{
		js.Global().Get("navigator").Get("geolocation").Call("getCurrentPosition",js.FuncOf(findStore))
	}else{
		js.Global().Get("alert").Invoke("Geolocation not supported by the browser!")
	}
	return nil
}


type Chat struct{}
var chat Chat

func (chat *Chat) startChat(this js.Value,inputs []js.Value) interface{}{
	
	js.Global().Get("document").Call("querySelector",".chat-popup").Get("style").Set("display","block")
	js.Global().Get("document").Call("querySelector",".send").Call("addEventListener","click",js.FuncOf(chat.sendMessage))

	return nil
}

func (chat *Chat) stopChat(this js.Value,inputs []js.Value) interface{}{

	js.Global().Get("document").Call("querySelector",".chat-popup").Get("style").Set("display","none")

	return nil
}

func (chat *Chat) sendMessage(this js.Value,inputs []js.Value) interface{}{

	msg:=js.Global().Get("document").Call("querySelector",".msg").Get("value").String()
	if msg!=""{
		clientText:=js.Global().Get("document").Call("createElement","div")
		clientText.Get("classList").Call("add","You")
		textString:=fmt.Sprintf("You:%v",msg)
		clientText.Set("innerHTML",textString)
		js.Global().Get("document").Call("querySelector",".textarea").Call("appendChild",clientText)
	}

	connection:=js.Global().Get("WebSocket").New("ws://localhost:8000/")
	
	connection.Set("onopen",js.FuncOf(func(this js.Value,inputs []js.Value) interface{}{
		connection.Call("send",msg)
		return nil
	}))

	connection.Set("onmessage",js.FuncOf(func(this js.Value,inputs []js.Value) interface{}{

		vibration:=js.Global().Get("navigator").Call("vibrate",10000)
		fmt.Println("vibrated:%v",vibration)
		serverText:=js.Global().Get("document").Call("createElement","div")
		serverText.Get("classList").Call("add","server")
		textString:=fmt.Sprintf("Server:%v",inputs[0].Get("data"))
		serverText.Set("innerHTML",textString)
		js.Global().Get("document").Call("querySelector",".textarea").Call("appendChild",serverText)
		return nil
	}))


	return nil
}


// MAIN FUNCTION


func main(){
	// c:=make(chan bool)
	document:=js.Global().Get("document")
	if js.Global().Get("Worker").Truthy()==true{
		// fmt.Println(js.Global().Get("Worker"))
		activityWorker=js.Global().Get("Worker").New("controllers/activityWorker.js")
		
		activityWorker.Set("onmessage",js.FuncOf(func(this js.Value,inputs []js.Value) interface{}{
			fmt.Println(inputs[0].Get("data").Get("count").String())
			activity:=inputs[0].Get("data").Get("desc")
			b:=[]byte(activity.String())
	  		json.Unmarshal(b, &logItems)
			fmt.Println(logItems)
			return nil
		}))

		document.Call("querySelector",".checkout").Call("addEventListener","click",js.FuncOf(func (this js.Value,inputs []js.Value) interface{}  {
			activityWorker.Call("terminate")
			js.Global().Get("alert").Invoke("Thankyou for shopping with us!")
			return nil
		}))
		

	}
	
	js.Global().Call("makeComponent", "product-element", js.FuncOf(customEl))
	js.Global().Call("makeComponent", "cart-element", js.FuncOf(customEl))

	if document.Get("readyState").String()!="loading" {
		respProducts:=product.getProducts()
		ui.displayProducts(respProducts)
		storage.createIndexDB(respProducts)
		
	}
	
	ui.setupCart();
	ui.displayViewedProducts()

	document.Call("querySelector",".search-box").Call("addEventListener","input",js.FuncOf(ui.filterFunction))

	document.Call("querySelector",".search-box").Call("addEventListener","keydown",js.FuncOf(
		func (this js.Value,inputs []js.Value) interface{} {
			keyCode:=inputs[0].Get("keyCode").String()
			if  strings.Contains(keyCode,"13"){
				search:=document.Call("querySelector",".search-box").Get("value").String()
				document.Call("querySelector",".search-box").Set("value","")
				ui.searchProducts(search)
				storage.setSessionStorage(search)
			}
			return nil
		}))

	// document.Call("addEventListener","DOMContentLoaded", js.FuncOf(func (js.Value,[]js.Value) interface{}  {
	// 	fmt.Println("Hello")
	// 	return nil
	// })) 

	respStores=store.getStores()
	document.Call("querySelector",".footer").Call("addEventListener","click",js.FuncOf(api.getLocation))

	document.Call("querySelector",".open-button").Call("addEventListener","click",js.FuncOf(chat.startChat))
	document.Call("querySelector",".cancel").Call("addEventListener","click",js.FuncOf(chat.stopChat))

	api.notifiyUsers()
	api.setupPushAPI()
	
	js.Global().Set("onload",js.FuncOf(func (js.Value,[]js.Value)interface{}  {
		fmt.Println("images are loaded")
		return nil
	}))
	// <-c
	select{}

}