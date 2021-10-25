// Cart structure
const cartTemplate= document.createElement('template');

cartTemplate.innerHTML=
`   <img  alt="Cartimage">
    <div>
        <h4></h4>
        <h5>€</h5>
        <span class="remove-item">Remove</span>
    </div>
    <style>
    .cart-item {
  display: grid;
  align-items: center;
  grid-template-columns: auto 1fr auto;
  grid-column-gap: 1.5rem;
  margin: 1.5rem 0;
}
img {
  width: 75px;
  height: 75px;
}
h4 {
  font-size: 0.85rem;
  text-transform: capitalize;
  letter-spacing: var(--mainSpacing);
}
h5 {
  margin: 0.5rem 0;
  letter-spacing: var(--mainSpacing);
}
.remove-item {
  color: grey;
  cursor: pointer;
}
    </style>
    
    `
;

class Cart extends HTMLElement{
    constructor(){    
        super();
        this.attachShadow({mode:'open'});
        this.shadowRoot.appendChild(cartTemplate.content.cloneNode(true));

        // Fetch the added item passed to web component
        let cartItem=JSON.parse(this.getAttribute("cartItem"));

        //set attributes of the struct of web component
        this.shadowRoot.querySelector("img").setAttribute("src",cartItem.image);
        this.shadowRoot.querySelector("h4").innerHTML=cartItem.title;
        this.shadowRoot.querySelector("h5").innerHTML=cartItem.price+"€";
        this.shadowRoot.querySelector("img").setAttribute("src",cartItem.image);
        this.shadowRoot.querySelector(".remove-item").setAttribute("data-id",cartItem.id);

        
    }

    // perform removal of selected node
    removeItem(id,price){

        activityWorker.postMessage("remove");

        let tempCart=localStorage.getItem("cart")?JSON.parse(localStorage.getItem("cart")):[]
        tempCart=tempCart.filter(item=>item.id!=id);
        localStorage.setItem("cart",JSON.stringify(tempCart));
        let self=this;
        self.remove();
        cartTotal=cartTotal-price;
        document.querySelector(".cart-total").innerHTML=cartTotal;
    }


    connectedCallback(){
        // Parse the cartItem passed and call event listener
        let cartItem=JSON.parse(this.getAttribute("cartItem"));
        this.shadowRoot.querySelector(".remove-item").addEventListener("click",()=>this.removeItem(cartItem.id,cartItem.price))
    }

    disconnectedCallback(){
        this.shadowRoot.querySelector(".remove-item").removeEventListener("click",()=>this.removeItem(cartItem.id,cartItem.price))

    }
}

window.customElements.define("cart-element", Cart)