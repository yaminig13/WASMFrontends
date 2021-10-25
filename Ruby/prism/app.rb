$cartArray=[]
$cartProducts=[]
$cartTotal=0
$recentlyViewedArray=[]
class  Cart< Prism::Component
    attr_accessor :title, :image, :price, :id 

    def initialize(cartHash)

        @title=cartHash["title"]
        @image=cartHash["image"]
        @price=cartHash["price"]
        @id=cartHash["id"]

    end

    def removeItem
        $cartProducts.each{|el|
            if el.title==title
                $cartProducts.delete(el)
            end   
        }

        $cartArray.each{|el|
            if(el['title']==title)
                $cartArray.delete(el)
            end
        }
        cartString=JSON.generate($cartArray)
        Storage.type("LocalStorage").set("cart",cartString)
        $cartTotal=$cartTotal-price


    end

    def render
        div(".cart-item",[
            img({
                props:{
                    src:@image,
                    alt:"CartImage"
                }
            }),
            div([
                h4({props:{innerHTML:@title}}),
                h5({props:{innerHTML:"#{@price}€"}}),
                span(".remove-item",{
                    onClick: call(:removeItem),
                    props:{innerHTML:"Remove"}
                })
            ])
        ])
    end
end


class Products< Prism::Component
    attr_accessor :title, :description, :image, :price, :id, :review, :productHash

    def initialize(productHash)
    
    puts "Initializing"
        @productHash=productHash
        @title=productHash['title']
        @description=productHash['description']
        @image=productHash['image']
        @price=productHash['price']
        @id=productHash['id']
        @review=productHash['review']
        @isBtnVisible="none"
    end


    def showBtn(val)
        if @isBtnVisible=="flex"
            @isBtnVisible="none"
        else
            @isBtnVisible=val
        end
    end

    def addCartItem
        Storage.type("LocalStorage").get("cart") do |response|
            if response.length()!=0
                $cartArray = JSON.parse(response)
            end
        end
        findItem= $cartArray.select{
            |obj| obj["title"]==title
        }
        if findItem.length()!=0
            Alert.showAlert("Already in cart")

        else
            Alert.showAlert("Successfully added to cart")
            
            $cartArray=$cartArray.push(@productHash)
            cartString=JSON.generate($cartArray)
            Storage.type("LocalStorage").set("cart",cartString)
        end
    end
    def navigate
        Storage.type("Cookie").get(nil) do |response|
            @ca=response.split(";")
        end
        cookieVal="#{title}=#{image};"

        Storage.type("Cookie").set(nil,cookieVal)

        if @ca.length()>2
            cookieVal="#{@ca[0]}=;expires=Thu, 01 Jan 1970 00:00:00 GMT"
            Storage.type("Cookie").set(nil,cookieVal)
        end

        $recentlyViewedArray=RecentlyViewedSection.new()

        Navigation.navigateto("http://localhost:8080/detail.html?id=",id)
    end
    def render
       article(".product",{
            onMouseOver: call(:showBtn).with("flex"),
            onMouseOut:call(:showBtn).with("flex")},[
                div(".img-container",[
                    img(".product-img",{
                        onClick:call(:navigate),
                        attrs:{
                            height:10,
                            width:10,
                        },
                        props:{
                            src:image
                        }
                    }),
                    button(".bag-btn",{
                        onClick: call(:addCartItem),
                        dataset:{
                            id:id
                        },
                        style:{
                            display:@isBtnVisible
                        }
                    },
                    [
                        i(".fas fa-shopping-cart","Add to Cart")
                    ])
                ]),
                h3({
                    props:{
                        innerHTML:title
                    }
                }),
                h4({
                    props:{
                        innerHTML:"#{price}€"
                    }
                }),
            ])
    end
end

class ChatMsg<Prism::Component

    def initialize(sender,msg)
        @msg=msg
        @sender=sender
    end

    def render
        div(".#{@sender}",{
            props:{innerHTML:"#{@sender}:#{@msg}"}
        })
    end
end

class DisplayWebSocket<Prism::Component
    attr_accessor :showChat

    def initialize
        @showChat=false
        @clientText=""
        @messageArray=[]
        @messageArray=@messageArray.push(ChatMsg.new("Server","Hello! How can we help you?"))

    end
    def startChat
        @showChat=true
    end 

    def stopChat
        @showChat=false
    end

    def sendMessage
        DOM.select(".msg").getValue do|response|
            @clientText=response
        end

        @messageArray=@messageArray.push(ChatMsg.new("You",@clientText))
        WebSocket.openConnection("ws://localhost:8000/",@clientText) do |response|
            VibrationAPI.vibrate(10000)
            @messageArray=@messageArray.push(ChatMsg.new("Server",response))
        end
    end

    def render
        div(".chatbot",[
            i(".fas fa-comment open-button",{
                onClick:call(:startChat)
            }),
            div(".chat-popup",{
                style:{
                        display:@showChat? 'block' : 'none'
                    }
            },[
                div(".form-container",[
                    h1("Chat"),
                    div(".textarea",
                        @messageArray
                    ),
                    div([
                        textarea(".msg",{
                            attr:{placeholder:"Type message.."}
                        })
                    ]),
                    button(".btn send",{
                        onClick:call(:sendMessage),
                        props:{innerHTML:"Send"}
                    }),
                    button(".btn cancel",{
                        onClick:call(:stopChat),
                        attr:{type:"button"},
                        props:{innerHTML:"Close"}
                    })
                ])
            ])
        ])
    end

end

class TrackLocation< Prism::Component

    def initialize
        $loadText="loaded"
        @stores=[]
        HTTP.get("http://localhost:8080/stores")do |response|
            @stores=JSON.parse(response.body)
        end
        @innerHTML="Find nearest Store"

         
    end

    def getLocation
        Location.find() do|response|
            findStore(response)
        end
    end

    def findStore(e)
        userLat=e['latitude']
        userLong=e['longitude']

        minDif=99999
        closest=nil
        for index in 0..@stores.length()-1
            lat1=userLat*(Math::PI/180)
            lon1=userLong*(Math::PI/180)
            lat2=@stores[index]['latitude']*(Math::PI/180)
            lon2=@stores[index]['longitude']*(Math::PI/180)
            R=6371
            
            x=(lon2-lon1)*Math.cos((lat1+lat2)/2)
            y=lat2-lat1
            dif=Math.sqrt(x*x+y*y)*R


            if dif<minDif
                closest=index
                minDif=dif
            end
        end
        @innerHTML=@stores[closest]["Address"]
    end
    
    def render
        div(".footer",{
            onClick:call(:getLocation)
            },[
            p(".location",{
                props:{innerHTML:@innerHTML}
            })
        ])
    end

end

class ActivityLog
    def initialize
        puts "hi"
        WorkerThread.startWorker('activityWorker.js') do|response|
            # response= JSON.parse(response)
            # puts response['data']
            puts response
        end

    end
end
class RecentlyViewedProducts<Prism::Component

    def initialize(path)
        @path=path
    end
    def render
    img({
        props:{
            src:"#{@path}"
        }
    })
    end
end
class RecentlyViewedSection<Prism::Component
    def initialize
        @showSection=false
        @viewedProducts=[]
        Storage.type("Cookie").get(nil) do |response|
            @ca=response.split(";")
        end

        if @ca.length()>0
            @showSection=true
        end

        for index in 0..@ca.length()-1
            valArray=@ca[index].split("=")
            val=valArray[1]
            @viewedProducts=@viewedProducts.push(RecentlyViewedProducts.new(val))
        end
    end

    def render
        section(".viewedProducts",{
            style:{display:@showSection? 'block' : 'none'},
        },[
            div(".section-title",[
                h2("Recently Viewed Products")
            ]),
            div(".viewed-product",
                @viewedProducts
            )
        ])
    end

end

class DisplayProducts < Prism::Component

    attr_accessor :cartVisibility


    def initialize
        @products=[]
        getProducts()
        puts "Initialization of Display products"
        @cartVisibility="none"
        @showCart=false
        @searchTerms=[]
        @searchHistory=nil
        @websocket=DisplayWebSocket.new()
        @location=TrackLocation.new()
        

    #    @activityLog=ActivityLog.new()
    # DOM.select('window').on('load') do |event|
    #   puts "hi"
    # end
        NotificationAPI.sendNotif("Good News","Sale on all products for limited time!")
        PushNotifAPI.sendPushNotif
    end

    def getProducts
        
        HTTP.get("http://localhost:8080/items") do |response|
            IndexedDB.createDB("Ruby",response.body)
            productsAJAX=JSON.parse(response.body)
            len=productsAJAX.length()-1
            (0..len).each do |i|
                @products=@products.push(Products.new(productsAJAX[i]))
            end
            # puts @products
        end
    end

    def cartVisibility
        Storage.type("LocalStorage").get("cart") do |response|
            if response.length()!=0
                $cartArray = JSON.parse(response)
            end
        end
        $cartProducts=[]
        $cartTotal=0
        $cartArray.each{|el|
            $cartProducts=$cartProducts.push(Cart.new(el))
            $cartTotal=$cartTotal+el['price']
        }

        @showCart=!@showCart
        # Navigation.getParam("id") do |response|
        #     puts response
        # end
    end

    def keydown(key,value)

        if key=="Enter"
            IndexedDB.openDB("Ruby",value) do|response|
                @searchProd=JSON.parse(response)
            
                puts @searchProd[0]
            
                if @searchProd.length()!=0
                    @products=[]
                    @products=@products.push(Products.new(@searchProd[0]))
        
                end
            end
            flag=false
            Storage.type("SessionStorage").get("searchHistory") do |response|
                @searchHistory=response
            end
            if @searchHistory!=nil
                @searchTerms=@searchHistory.split(",")

                if @searchTerms.include?(value)
                    flag=true
                end
                if flag==false
                    @searchHistory="#{@searchHistory},#{value}"
                end
            else
                @searchHistory=value
            end
            Storage.type("SessionStorage").set("searchHistory",@searchHistory)

        end
    end

    def filterFunction(value)
        Storage.type("SessionStorage").get("searchHistory") do |response|
                @searchHistory=response
            end
        if @searchHistory!=nil
            result=nil
            sessionTerms=[]
            filter=value.downcase

            sessionTerms=@searchHistory.split(",")

            for term in sessionTerms
                if term.index(value)!=nil
                    result="#{result} #{term}"
                end
            end
        end
    end

    def activityLog
        WorkerThread.startWorker('activityWorker.js') do|response|
            # response= JSON.parse(response)
            # puts response['data']
            WorkerThread.receiveMsg(response) do |msg|
                puts msg
            end
        end
        $recentlyViewedArray=RecentlyViewedSection.new()
    end
    
    def render
        body([
            nav(".navbar",[
                img({
                    onLoad:call(:activityLog),
                    style:{
                        display:'none'
                    },
                    props:{src:'/assets/images/airfryer.jpg'}
                }),
                div(".navbar-center",[
                    div(".search",{
                        onKeydown: call(:keydown).with_event_data(:key).with_target_data(:value).stop_propagation
                    },
                    [
                        input(".search-box",
                            onInput: call(:filterFunction).with_target_data(:value).stop_propagation,
                            attr:{
                                type:"text",
                                # list:$productSuggestion
                            }
                        ),
                        # datalist("#productSuggestion"),
                        span(".nav-icon",[
                            i(".fa fa-search")
                        ])
                    ]),
                    div([
                        h1("Elektronika")
                    ]),
                    div(".cart-btn",{
                        onClick: call(:cartVisibility)
                    },
                    [
                        span(".nav-icon",[
                            i(".fas fa-cart-plus")
                        ])
                    ])
                ])
            ]),
            section(".products",[
                div(".section-title",[
                    h2("Our Products")
                ]),
                div(".products-center",
                @products 
                )
            ]),
            div(".cart-overlay",{
                class:{
                    transparentBcg:@showCart
                }
            },
            [
                div(".cart",{
                    style:{
                        display:@showCart? 'inline' : 'none'
                    }
                },
                [
                    span(".close-cart",{
                        onClick: call(:cartVisibility)
                    },
                    [
                        i(".fas fa-window-close")
                    ]),
                    h2("Your Cart"),
                    div(".cart-content",
                    $cartProducts,
                    {
                        props:{innerHTML:""}
                        
                    }),
                    div(".cart-footer",[
                        h3({props:{innerHTML:"Total Amount: €"}},[
                            span(".cart-total",{props:{innerHTML:$cartTotal}})
                        ])
                    ])
                ])
            ]),
            @websocket,
            @location,
            $recentlyViewedArray
        ])
        
    end
end


Prism.mount(DisplayProducts.new)
