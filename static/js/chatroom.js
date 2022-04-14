$(document).ready(function () {
    listWs = new WebSocket("ws://192.168.28.4:5050/chatroom/userlist")
    listWs.onopen=function(){
        console.log("connected");
    }
    listWs.onmessage = function(e){
        console.log(e);
    }
})