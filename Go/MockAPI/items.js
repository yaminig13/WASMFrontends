const { string, boolean } = require('casual');
const casual = require('casual')

// module.exports = () => {
//     casual.define('item', function() {
//         var name=["Electric Mixer","Rice Cooker","Hair Dryer","Stand Mixer","Refrigerator","Microwave","Oven","Induction","Hair Straightner"];
//         return {
//             title: casual.random_element(name),
//             description: casual.string,
//             path:String,
//             price: casual.double(from=100,to=1000),
//             available: Boolean,
//             quantity: casual.integer(from=1,to=100),
//             id: casual.uuid,
//         }
//     })

//     const data = {
//         items: [],
//     }
    
//     data.items.push(casual.item)
//     // Create 100 users
//     for (let i = 0; i < 100; i++) {
//         tempItem=casual.item;
//         flag=false;
//         data.items.forEach(element => {
//             if(element.title==tempItem.title){
//                 flag=true;
//             }     
//         }); 
//         if(!flag)
//             data.items.push(tempItem);  
//     }
//     return data
// }

module.exports=()=>{

    let items=[
        {
            "title": "oven",
            "description": "Built in design for eye level in a kitchen column, Large 74L fan assisted main oven and grill, 37L conventional oven and grill, Dimensions-(H)89cm x (W)60cm x (D)58cm, Enamel lined oven interiors for easier cleaning",
            "image":"assets/images/oven.jpg",
            "price": casual.integer(from=100,to=1000),
            "id": "af0105b8-8ee1-4056-9b33-970d4d74d830",
            "review":[{
                "username":"Theresa",
                "image":String,
                "description":"I found the oven super useful!!!",
                "date":Date
            }]
        },

        {
            "title": "refrigerator",
            "description": "Sunsang - 28 Cu. Ft. French Door Refrigerator with CoolSelect Pantry™ - Fingerprint Resistant Stainless Steel",
            "image":"assets/images/refrigerator.jpg",
            "price": casual.integer(from=100,to=1000),
            "id": "5c997c5e-0e5b-48a2-85ea-48eb0ddd574b",
            "review":[{
                "username":String,
                "image":String,
                "description":String,
                "date":Date
            }]
        },
        {
            "title": "headset",
            "description": "Rose - QuietComfort 35 II Gaming Headset – Comfortable Noise Cancelling Headphones - Black",
            "image":"assets/images/headset.jpg",
            "price": casual.integer(from=100,to=1000),
            "id": "d2582006-8a84-44ec-a485-437876b8798e",
            "review":[{
                "username":String,
                "image":String,
                "description":String,
                "date":Date
            }]
        },
        {
            "title": "airfryer",
            "description": "Balle Pro Series - 2-qt. Touchscreen Air Fryer - Black Matte",
            "image":"assets/images/airfryer.jpg",
            "price": casual.integer(from=100,to=1000),
            "id": "c05d297a-f98a-44cf-9c65-3e72303f57ad",
            "review":[{
                "username":String,
                "image":String,
                "description":String,
                "date":Date
            }]
        },
        {
            "title": "speaker",
            "description": "KCM CHARGE waterproof speaker",
            "image":"assets/images/speaker.jpeg",
            "price": casual.integer(from=100,to=1000),
            "id": "bc1a1f4a-a2b3-4a20-bd51-005855e03b6d",
            "review":[{
                "username":String,
                "image":String,
                "description":String,
                "date":Date
            }]
        }

    ]

    

    return items

}