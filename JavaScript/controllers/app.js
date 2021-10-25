
const productsDOM=document.querySelector(".products-center");
const cartDOM=document.querySelector(".cart");
const cartOverlay=document.querySelector(".cart-overlay");
var searchHistory="";
var cart=[];
var cartTotal=null;
var stores=[];
var connection;
var ui;
var activityWorker;

function urlBase64ToUint8Array(vapidPublicKey){
    const padding = "=".repeat((4 - vapidPublicKey.length % 4) % 4);
                    const base64 = (vapidPublicKey + padding)
                    .replace(/\-/g, "+")
                    .replace(/_/g, "/");

                    const rawData = window.atob(base64);
                    const outputArray = new Uint8Array(rawData.length);

                    for (let i = 0; i < rawData.length; ++i) {
                        outputArray[i] = rawData.charCodeAt(i);
                    }
                    return outputArray;
}

class DataStore{
    async getProducts(){
        let response= await fetch("http://localhost:8080/items");
        let products= await response.json();        
        return products;
    }

    async getStores(){
        let response= await fetch("http://localhost:8080/stores");
        stores= await response.json();        
        return stores;
    }
}

class UI{    

    // search product from indexdb
    searchProducts(searchTerm){
        let self=this;
        searchTerm=searchTerm.toLowerCase();
        const dbName="JavaScript";
        var request=indexedDB.open(dbName,1);
        request.onerror=function(event){
            console.error("error found");
        }

        request.onsuccess=function(event){
            console.log("connection open");
            var db=event.target.result;
            var txn=db.transaction("products","readwrite");   
            let tempProducts=txn.objectStore("products");
            let titleIndex=tempProducts.index("title")
            let query=titleIndex.getAll(searchTerm);

            query.onsuccess=function(){
                if(query.result!==undefined){
                    self.displayProducts(query.result);
                }
                else{
                    console.log("no product found")
                }
            }
        }
    }

    // display fetched product
    displayProducts(products){
        let result="";
        products.forEach(product=>{
            product.description="";
            product.review=null;            
            result+=`
            <product-element product=${JSON.stringify(product)} ></product-element>
            `
        });
        productsDOM.innerHTML=result;       
    }

    // search box input filter
    filterFunction(e){
        
        if(sessionStorage.getItem("searchHistory")){
            let result="";
            let sessionTerms=[];
            let filter=e.target.value.toLowerCase();
            sessionTerms=sessionStorage.getItem("searchHistory").split(",");
            sessionTerms.forEach(term=>{
                if(term.indexOf(filter)>-1){
                    result+=`
                    <option value="${term}">
                    `
                }
            });
            document.querySelector("#productSuggestion").innerHTML=result;
        }
    }

    
    setupCart(){
        var cartBtn=document.querySelector(".cart-btn");
        var closeCartBtn=document.querySelector(".close-cart");
        cartBtn.addEventListener("click",()=>this.showCart());
        closeCartBtn.addEventListener("click",()=>this.hideCart());
    }

    // initialize cart
    populateCart(){
        cart=[];
        cartTotal=0;
        document.querySelector(".cart-content").innerHTML="";
        document.querySelector(".cart-total").innerHTML=cartTotal;
        if(typeof(Storage)!=="undefined"){

            cart=localStorage.getItem("cart")?JSON.parse(localStorage.getItem("cart")):[]
            
        }
        cart.forEach(cartItem=>{
            var div=document.createElement('div');
            div.classList.add("cart-item");        
            div.innerHTML=`<cart-element cartItem=${JSON.stringify(cartItem)}></cart-element>`
            cartTotal+=cartItem.price;
            document.querySelector(".cart-content").appendChild(div);
        })
        document.querySelector(".cart-total").innerHTML=cartTotal;

    }
    
    // fill cart before showing
    showCart(){
        this.populateCart();
        cartOverlay.classList.add("transparentBcg");
        cartDOM.style.display="inline";
    }
    hideCart(){
        cartOverlay.classList.remove("transparentBcg");
        cartDOM.style.display="none";
    }

    displayViewedProducts(){
        let ca=document.cookie.split(";");
        if(ca.length){
            document.querySelector(".viewedProducts").style.display="block";
        }
        for(let i=0;i<ca.length;i++){
            let val=ca[i].split("=")[1];
            let imgDiv=document.createElement("img");
            // let val=getCookie(ca[i])
            imgDiv.setAttribute("src",val);
            document.querySelector(".viewed-product").appendChild(imgDiv);
        }
    }
}

class Storage{

    createIndexDB(products){
        const dbName="JavaScript";
        var request=indexedDB.open(dbName,1);
        request.onerror=function(event){
            console.error("error found");
        }

        request.onupgradeneeded=function(event){
            var db=event.target.result;

            var objectStore=db.createObjectStore("products",{keyPath:"id"});
            objectStore.createIndex("title","title",{unique:false});
            objectStore.transaction.oncomplete=function(event){
                var productStore=db.transaction("products","readwrite").objectStore("products");
                products.forEach(function(product){
                    productStore.add(product);
                })
            }   
        }
    }

    setSessionStorage(searchTerm){
        let flag=false;
        if(typeof(Storage)!=="undefined"){
            if(sessionStorage.getItem("searchHistory")){
                searchHistory=sessionStorage.getItem("searchHistory");
                let sessionTerms=sessionStorage.getItem("searchHistory").split(",");

                // check if term already added
                sessionTerms.forEach(term=>{
                    if( term==searchTerm){
                        flag=true;
                    } 
                })
                if(!flag){
                    searchHistory+=","+searchTerm;
                }
            }
            else{
                searchHistory=searchTerm;
            }
            sessionStorage.setItem("searchHistory",searchHistory);
        }
    }
}

class API{
    getLocation(stores){
        if(navigator.geolocation){
            navigator.geolocation.getCurrentPosition(this.findStore);
        }
        else{
            alert("GeoLocation not supported by the browser");
        }
        
        	
        
    }

    findStore(position){
        console.log(stores)
        var userLat=position.coords.latitude;
        var userLong=position.coords.longitude;

        var minDif=99999;
        var closest;

        for(let index=0;index < stores.length;++index){
                var lat1=userLat*Math.PI/180;
                var lon1=userLong*(Math.PI/180);
                var lat2=stores[index].latitude*Math.PI/180;
                var lon2=stores[index].longitude*Math.PI/180;

                var R=6371;
                var x=(lon2-lon1)*Math.cos((lat1+lat2)/2);
                var y=(lat2-lat1);
                var dif=Math.sqrt(x*x+y*y)*R;
            
            console.log(dif)
            if(dif<minDif){
                closest=index;
                minDif=dif;
            }

        }
        console.log(closest)
        document.querySelector(".location").innerHTML=stores[closest].Address;
    
    }

    notifyUsers(){
        if(window.Notification && Notification.permission!=="denied"){
            Notification.requestPermission(function(status){
                var notif= new Notification("Good News",{
                    body:"Sale on all products for limited time!"
                })
            })
        }
    }

    setupPushAPI(){
        navigator.serviceWorker.register("serviceWorker.js");

        navigator.serviceWorker.ready.then(function(registration){
            return registration.pushManager.getSubscription().then(async function(subscription){
                if(subscription){
                    return subscription;
                }

                const response=await fetch("http://localhost:8080/vapidPublicKey");
                var vapidPublicKey=await response.text();
                const convertedVapidKey=urlBase64ToUint8Array(vapidPublicKey);
                return registration.pushManager.subscribe({
                    userVisibleOnly:true,
                    applicationServerKey:convertedVapidKey
                });
            })
        }).then(function(subscription){
            fetch("http://localhost:8080/register",{
                method:"post",
                headers:{
                    "Content-type":"application/json"
                },
                body:JSON.stringify({
                    subscription:subscription
                }),
            });
            
            fetch("http://localhost:8080/sendNotification",{
                method:"post",
                headers:{
                    "Content-type":"application/json"
                },
                body:JSON.stringify({
                    subscription:subscription,
                    delay:120
                }),
            });
        })
    }
}

class Chat{

    startChat(){
        document.querySelector(".chat-popup").style.display="block";
        document.querySelector(".send").addEventListener("click",()=>this.sendMessage());
    }

    stopChat(){
        document.querySelector(".chat-popup").style.display="none";
    }

    sendMessage(){
        let msg=document.querySelector(".msg").value;
        if(msg){
            var clientText=document.createElement("div");
            clientText.classList.add("You");
            clientText.innerHTML="You:"+msg
            document.querySelector(".textarea").appendChild(clientText);
        }
        connection=new WebSocket("ws://localhost:8000/");


        connection.onopen=()=>{
            
            connection.send(msg);
        }

        connection.onmessage=(e)=>{
            var vibration=window.navigator.vibrate(10000);
            console.log("vibrated:"+vibration);
            var serverText=document.createElement("div");
            serverText.classList.add("server");
            serverText.innerHTML="Server:"+e.data
            document.querySelector(".textarea").appendChild(serverText);
        }        
    }
}

document.addEventListener("DOMContentLoaded",()=>{

    var products =new DataStore();
    var ui= new UI();
    var storage=new Storage();
    var api=new API();
    var store=new DataStore();
    var chat=new Chat();


    // call to fetch products and display them and store in indexdb
    api.notifyUsers();
    products.getProducts().then(products=>{
        ui.displayProducts(products);
        storage.createIndexDB(products);
    })

    // add basic functionalities to cart

    ui.setupCart();

    // search functionality, listen for new input in searchbox
    document.querySelector(".search-box").addEventListener("input",ui.filterFunction);

    // search on enter
    document.querySelector(".search-box").addEventListener("keydown",function(event){
        if(event.keyCode==13){
            let search=document.querySelector(".search-box").value;
            document.querySelector(".search-box").value="";
            ui.searchProducts(search)
            storage.setSessionStorage(search);
        }
    })
    
    store.getStores().then(stores=>{
        document.querySelector(".footer").addEventListener("click", ()=>api.getLocation());

    })

    document.querySelector(".open-button").addEventListener("click",()=>chat.startChat());
    document.querySelector(".cancel").addEventListener("click",()=>chat.stopChat());
    api.setupPushAPI();
    ui.displayViewedProducts();

if (window.Worker) {
    activityWorker=new Worker("controllers/activityWorker.js");
    activityWorker.onmessage=function(e){
        console.log("")
        console.log("Number of times clicked:"+e.data.count);
        e.data.desc.forEach(element => {
            console.log(element.type+":"+element.timeStamp);
        });
        

    }
    document.querySelector(".checkout").addEventListener("click",function(){
        activityWorker.terminate();
        alert("Thankyou for shopping with us")
    })
}
});

// wait for image to load
window.onload= function(event){
    console.log("images have been loaded too")
}



