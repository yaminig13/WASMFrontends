package main

import (
	"gitlab.com/microo8/golymer"
)

var template=golymer.NewTemplate(`<article class="product">
    <div class="img-container">
        <img height=10 width=10 class="product-img">
        <button class="bag-btn" >
            <i class="fas fa-shopping-cart"></i>
            Add to cart
        </button>
    </div>
    <h3></h3>
    <h4></h4>
</article>

<style>
.img-container {
  position: relative;
  overflow: hidden;
}
.bag-btn {
  position: absolute;
  top: 70%;
  right: 0;
  background: var(--primaryColor);
  border: none;
  text-transform: uppercase;
  padding: 0.5rem 0.75rem;
  letter-spacing: var(--mainSpacing);
  font-weight: bold;
  cursor: pointer;
  display:none
}
.bag-btn:hover {
  color: var(--mainWhite);
}
.fa-shopping-cart {
  margin-right: 0.5rem;
}
.product-img {
  display: block;
  width: 100%;
  min-height: 10rem;
}

.product h3 {
  text-transform: capitalize;
  font-size: 1.1rem;
  margin-top: 1rem;
  letter-spacing: var(--mainSpacing);
  text-align: center;
}

.product h4 {
  margin-top: 0.7rem;
  letter-spacing: var(--mainSpacing);
  color: var(--primaryColor);
  text-align: center;
}
</style>`)

type ProdElement struct{
	golymer.Element
	// dom.Element
	
}

func NewElement() *ProdElement {
	
	e:=&ProdElement{}
	// fmt.Println(e.Element.AttachShadow{}
	e.Element.SetTemplate(template)
	// e.this.Set(e.this,"super")
	// opts:=dom.AttachShadowOpts{
	// 	Open:true,
	// }
	// // fmt.Println(e.this)
	// shadowRoot:=e.Element.AttachShadow(opts)
	// fmt.Println(shadowRoot)
	// opts.open=true;
	// dom.AttachShadow(opts)
		return e

}




func main() {
	err:=golymer.Define(NewElement)
	if err!=nil{
		panic(err)
	}
	// p:=dom.Element{}
	// template := js.Global().Get("document").Call("createElement","template")

// 	templateHTML:=`<article class="product">
//     <div class="img-container">
//         <img height=10 width=10 class="product-img">
//         <button class="bag-btn" >
//             <i class="fas fa-shopping-cart"></i>
//             Add to cart
//         </button>
//     </div>
//     <h3></h3>
//     <h4></h4>
// </article>

// <style>
// .img-container {
//   position: relative;
//   overflow: hidden;
// }
// .bag-btn {
//   position: absolute;
//   top: 70%;
//   right: 0;
//   background: var(--primaryColor);
//   border: none;
//   text-transform: uppercase;
//   padding: 0.5rem 0.75rem;
//   letter-spacing: var(--mainSpacing);
//   font-weight: bold;
//   cursor: pointer;
//   display:none
// }
// .bag-btn:hover {
//   color: var(--mainWhite);
// }
// .fa-shopping-cart {
//   margin-right: 0.5rem;
// }
// .product-img {
//   display: block;
//   width: 100%;
//   min-height: 10rem;  
// }

// .product h3 {
//   text-transform: capitalize;
//   font-size: 1.1rem;
//   margin-top: 1rem;
//   letter-spacing: var(--mainSpacing);
//   text-align: center;
// }

// .product h4 {
//   margin-top: 0.7rem;
//   letter-spacing: var(--mainSpacing);
//   color: var(--primaryColor);
//   text-align: center;
// }
// </style>`
	// template.Set("innerHTML",templateHTML)

	// js.Global().Get("customElements").Call("define","product-element",js.FuncOf(NewElement))
select{}
}