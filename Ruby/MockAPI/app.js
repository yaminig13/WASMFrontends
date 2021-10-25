let express = require('express');
let fs = require('fs');
let path = require('path');
const WebSocket = require('ws');
let webpush = require('web-push');
var multer  = require('multer');
let app = express();
var atob=require('atob');

app.use(express.json({limit:'50mb'}));
app.use(express.urlencoded({ limit:'50mb', extended: true }));
// app.use(express.static(path.join(__dirname,'../JavaScript')))
app.use(express.static(path.join(__dirname,'../prism')))


let items= require('./items')
let stores= require('./stores')

var upload = multer({ dest: '/uploads/' });

// configure the app to use bodyParser()
// app.use(bodyParser.urlencoded({
//     extended: true
// }));

// app.use(function (req, res, next) { setTimeout(next, 400) });

app.use(function (req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    // res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    next();
});

// app.get("/detail.html",function(request,response){

//     console.log(request.query.id)
//     let res=items().find(item=>item.id=request.params.id)
//     response.status(200).jsonp(res)

// })
app.post('/review', (request, response) => {
    if (request.method === 'POST') {
        var image64=request.body.image64;
        image64=image64.replace('data:image/png;base64','');
        image64=image64.replace(' ','+');
        response.status(200).send();
        // upload.single(atob(image64));
//         fs.writeFile("/uploads"+"/out.png", image64, 'base64', function(err) {
//             console.log(err);
//     });
// fs.readFileSync(image64, {encoding: 'base64'});
    }
})
app.get('/items/:id', (request, response) => {
    if (request.method === 'GET') {
    	let res=items().find(item=>item.id==request.params.id)
        response.status(200).jsonp(res)
    }
})

app.get('/items', (request, response) => {
    if (request.method === 'GET') {
        console.log("Items called")
        response.status(200).jsonp(items())
    }
})

app.get('/stores', (request, response) => {
    if (request.method === 'GET') {
        response.status(200).jsonp(stores())
    }
})

// Websocket server setup

const wss = new WebSocket.Server({ port: 8000 })
 
wss.on('connection', ws => {
    ws.on('message', message => {    
        if(message=="how do you deliver products?" ){
            ws.send("we have only click and collect option available as of now.")        
        }
        else{
            ws.send("Please contact our customer service. ")
        }    
    })
})

// push notification backend
const publicKey = "BBCaPWiiDVFRL746EBAlsLMe428QKD6IvDXyKz9RerCYffVP_G8zHUSj8Dh1wJI55e3K1EWVz03IuUKwEdBYMvA";
const privateKey ="dhLcOG4NuLVCKXr7T5xdjlmDygfKpAtrDRN5cbifrWk"
webpush.setVapidDetails('mailto:helloyamini@gmail.com',publicKey,privateKey);

// subscribe route
app.get('/vapidPublicKey', function(req, res) {
    res.send(publicKey);
  });

app.post('/register',(req,res) => {
        res.sendStatus(201);

    // const subscription = req.body;
    // res.status(201).json({});
    // const payload=JSON.stringify({title:'push test'})
    // webpush.sendNotification(subscription,payload).catch(err=>console.error(err))
})


app.post('/sendNotification', function(req, res) {
    const subscription = req.body.subscription;
    const payload = null;
    const options = {
      TTL: req.body.ttl
    };

    setTimeout(function() {
      webpush.sendNotification(subscription, payload, options)
      .then(function() {
        res.sendStatus(201);
      })
      .catch(function(error) {
        res.sendStatus(500);
        console.log(error);
      });
    }, req.body.delay * 1000);
  });



app.listen(8080,()=>{
    console.log('JSON Server is running at Port 8080');
});