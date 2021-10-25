var clicks={
    count:0,
    desc:[]
};

var removeClicks=0;
onmessage=function(e){
    if(e.data=="add"){

        clicks.count+=1;
        let click={
            type:"Add",
            timeStamp:new Date()
        }
        clicks.desc.push(click)
    }
    if(e.data=="remove"){
        clicks.count+=1;
        let click={
            type:"Remove",
            timeStamp:new Date()
        }
        clicks.desc.push(click)
    }
}
setInterval(function(){
        postMessage(clicks)},30000);