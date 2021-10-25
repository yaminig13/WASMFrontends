class ReviewForm< Prism::Component

    def initialize
        @showCanvas=false
        @showBtn=true
    end
    def render
        div([
            h1("Add Review"),
            div(".user",[
                label("Name"),
                input()
            ]),
            div(".reviewText",[
                label("Review"),
                textarea(),
                input(".cameraInput",{
                    attrs:{
                        name:"cameraInput",
                        type:"file",
                        accept:"image/png,image/jpeg,image/jpg"
                    },
                    style:{
                        visibility:'hidden'
                    }
                })
            ]),
            div(".cameraCanvas",{
                style:{
                        display:@showCanvas? 'block' : 'none'
                    }},
                [video("#video",{
                    attrs:{
                        width:200,
                        height:200
                    }
                }),
                canvas("#canvas",{
                    attrs:{
                        width:200,
                        height:200
                    }
                }),
                div([
                    button(".snap",{
                        attrs:{type:"button"},
                        props:{innerHTML:"Snap"}
                    })
                ])
            ]),
            div(".reviewBtns",[
                    button(".reviewImagebtn",
                        {
                            onClick:call(:startCamera),
                            style:{display:@showBtn? 'block' : 'none'},
                            props:{innerHTML:"Take Picture"},
                            attrs:{type:"button"}
                            
                        }),
                    button(".reviewSubmit",{
                        onClick:call(:submitReview),
                        props:{innerHTML:"Submit Review"},
                        attrs:{type:"button"}
                    })
                ])
        ])
    end

    def submitReview
        HTTP.post("http://localhost:8080/review/",@imageData)
    end

    def startCamera
        @showCanvas=true
        @showBtn=false

        CameraAPI.enable_camera do |response|
            @imageData=response
        end

        # DOM.select("canvas").getElement do |response|
        #     puts response
        # end
    end
end

class MainStruct < Prism::Component

    def initialize
        @productJSON={}
        @reviewForm=ReviewForm.new()
        

    end

    def getid
        Navigation.getParam("id") do|response|
            @id=response
        end  
      
        HTTP.get("http://localhost:8080/items/#{@id}") do |response|
            @productJSON=JSON.parse(response.body)
        end
    end

    def detail
        div(".detail",[
            img({
                    onLoad:call(:getid),
                    style:{
                        display:'none'
                    },
                    props:{src:'/assets/images/airfryer.jpg'}
                }),
            div(".card",[
                img({
                    props:{src:@productJSON['image']}
                })
            ]),
            div(".info",[
                h1(".title",{
                    props:{innerHTML:@productJSON['title']}
                }),
                p(".price",{
                    props:{innerHTML:"#{@productJSON['price']}€"}
                }),
                p(".description",{
                    props:{innerHTML:@productJSON['description']}
                })
            ])
        ])
    end

    def render
        body(".main",[
            div(".navbar",{
                onClick:call(:navigateid)
            },[
                div(".navbar-center",[
                    div([
                        h1("Elektronika")
                    ])
                ])
            ]),
            detail,
            @reviewForm
        ])
    end
end

Prism.mount(MainStruct.new)
