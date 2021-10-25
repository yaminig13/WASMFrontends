async function getProducts(id){
        let response= await fetch("http://localhost:8080/items/"+id);
        let product= await response.json();        
        return product;
    }

function displayData(product){
    document.querySelector("img").src=product.image;
    document.querySelector(".title").innerHTML=product.title;
    document.querySelector(".price").innerHTML=product.price+"€";
    document.querySelector(".description").innerHTML=product.description;
}

// function loadReviews(product){
//         let revieWrapper=document.querySelector(".reviewItems");
//         product.review.forEach(element => {
//             let result="";
//             let reviewDiv=document.createElement("div");
//             reviewDiv.classList.add("review");
//             result=`<h4>${element.username}</h4>
//             <p>${element.description}</p>
//             <img src="${element.image}">
//             `
//             reviewDiv.innerHTML=result;
//             revieWrapper.appendChild(reviewDiv);
//         });

// }
// function dataURLtoBlob(dataURL) {
//   let array, binary, i, len;
//   binary = atob(dataURL.split(',')[1]);
//   array = [];
//   i = 0;
//   len = binary.length;
//   while (i < len) {
//     array.push(binary.charCodeAt(i));
//     i++;
//   }
//   return new Blob([new Uint8Array(array)], {
//     type: 'image/png'
//   });
// };
function startCamera(){
    document.querySelector(".cameraCanvas").style.display="block";
        document.querySelector(".reviewImagebtn").style.display="none";

    var video = document.getElementById("video");
    var canvas=document.getElementById("canvas");
    var context=canvas.getContext("2d");

    document.getElementById("snap").addEventListener("click",function(){
        context.drawImage(video,0,0,200,200);
        var imageData=canvas.toDataURL();
        document.querySelector(".cameraInput").setAttribute("value",imageData);

//                 const file = dataURLtoBlob( canvas.toDataURL() );

//         const fd = new FormData;

// fd.append('image', file);
        document.querySelector(".reviewSubmit").addEventListener("click",()=>{
            fetch("http://localhost:8080/review/",{
                method:"POST",
                headers:{
                    'content-type': 'application/json; charset=utf-8',                
                },
                body:JSON.stringify(
                {
                    image64:imageData
                }),
            }).then(()=>{alert("Review Submitted!")})

        //     // });
        //     fetch(imageData)
        //     .then(res => res.blob())
        //     .then(blob => {
        //         const fd = new FormData();
        //         const file = new File([blob], "filename.jpeg");
        //         fd.append('image', file)
  
            
        //         const API_URL = 'http://localhost:8080/review/'
        //         fetch(API_URL, {method: 'POST', body: fd})
        //         .then(res => res.json()) 
        //         .then(res => console.log(res))
        //     });
        })
        // // console.log(imageData)
        //     // canvas.toBlob(res,"image/jpeg");
            
    })
    if(navigator.mediaDevices && navigator.mediaDevices.getUserMedia){
        navigator.mediaDevices.getUserMedia({video:true}).then(function(stream){
            video.srcObject=stream;
            video.play();
        })
    }
}


var url_string=window.location.href;
var url=new URL(url_string);
var id=url.searchParams.get("id");
console.log(id)
getProducts(id).then(product=>{
    displayData(product)
    document.querySelector(".reviewImagebtn").addEventListener("click",()=>startCamera());

    // loadReviews(product)
});                         