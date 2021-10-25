// Product structure
const template= document.createElement('template');

template.innerHTML=`
<article class="product">
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
</style>
`

class Product extends HTMLElement{
    constructor(){
        super();
        this.attachShadow({mode:'open'});
        this.shadowRoot.appendChild (template.content.cloneNode(true));

        // Fetch the added item passed to web component

        let customProduct=JSON.parse(this.getAttribute("product"));

        //set attributes of the struct of web component

        this.shadowRoot.querySelector("img").setAttribute("src",customProduct.image);
        this.shadowRoot.querySelector("button").setAttribute("data-id",customProduct.id);
        this.shadowRoot.querySelector("h3").innerHTML=customProduct.title;
        this.shadowRoot.querySelector("h4").innerHTML=customProduct.price+"€";

    }

    // mouseover hide and show buttons
    hideBtn(){
        this.shadowRoot.querySelector(".bag-btn").style.display="none";
    }

    showBtn(){
         this.shadowRoot.querySelector(".bag-btn").style.display="flex";
    }

    // adding the product to cart
    getCartBtns(){
        let customProduct=JSON.parse(this.getAttribute("product"));
        let button=this.shadowRoot.querySelector(".bag-btn");
       
        button.addEventListener('click',()=>{

          activityWorker.postMessage("add");

            // fetch products already in cart
            let localCart=localStorage.getItem('cart')?JSON.parse(localStorage.getItem("cart")):[];

            // find if product passed is already in cart
            let inCart=localCart.find(item=>item.id===customProduct.id);
            if(inCart){
                alert("Already in cart!");
            }
            else{
                localCart.push(customProduct);
                if(typeof(Storage)!=="undefined"){
                    localStorage.setItem("cart",JSON.stringify(localCart));
                } 
                alert("Successfully added to cart!")
            }
        })
    }

    navigate(e){      
        let customProduct=JSON.parse(this.getAttribute("product"));
        window.location=("http://localhost:8080/detail.html?id="+customProduct.id);

        let ca=document.cookie.split(";");
        document.cookie=customProduct.title+"="+customProduct.image+";"
        if(ca.length>2){
          document.cookie=ca[0]+"=;"+"expires=Thu, 01 Jan 1970 00:00:00 GMT";
        }
        ui.displayViewedProducts();
    }

    connectedCallback(){  
        
        // add functionalities to product button
        this.getCartBtns();
        this.shadowRoot.querySelector(".product").addEventListener("mouseover",()=>this.showBtn());
        this.shadowRoot.querySelector(".product").addEventListener("mouseleave",()=>this.hideBtn());
        this.shadowRoot.querySelector("img").addEventListener("click",()=>this.navigate());
    }
    disconnectedCallback(){
        this.shadowRoot.querySelector(".product").removeEventListener("mouseover",()=>this.showBtn(),true);
        this.shadowRoot.querySelector(".product").removeEventListener("mouseleave",()=>this.hideBtn(),true);

    }
}

window.customElements.define("product-element", Product)