self.addEventListener("push",function(event){
    event.waitUntil(
        self.registration.showNotification("Hi There!",{
         body:"Don't forget to leave your feedback! Happy Shopping!",   
        })
    );
});