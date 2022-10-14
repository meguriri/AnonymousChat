function getcurTime() {
    var Digital=new Date();
    var hours=Digital.getHours();
    var minutes=Digital.getMinutes();
    var seconds=Digital.getSeconds();
    if(minutes<=9){
        minutes="0"+minutes;
    }if(seconds<=9){
        seconds="0"+seconds;
    }
    myclock=hours+":"+minutes+":"+seconds;
    return myclock;
}

$(document).ready(function () {
    listWs = new WebSocket("ws://localhost:5050/chatroom/userlist")
    listWs.onopen=function(){
        console.log("userlist connected");
    }


    listWs.onmessage = function(e){
        let user=JSON.parse(e.data)
        console.log(user)
        let c
        $('#listinfo').html("")
        for(i=0;i<user.length;i++){
            console.log(user[i])
            if(user[i].gender==0){
                c=' ♂ '
            }else{
                c=' ♀ '
            }
            $('#listinfo').append(
                '<li class="my-1 border-0 list-group-item d-flex p-1">\n' +
                '        <span class="p-1 rounded-circle text-white text-center" style="width: 40px;height: 40px;' +
                'background-color: rgb('+user[i].color[0]+','+user[i].color[1]+','+user[i].color[2]+')"><i\n' +
                '            class="bi bi-person"></i> </span>\n' +
                '        <p class="mx-1 text-black-50">'+user[i].nickname+c+'️</p>\n' +
                '    </li>'
            )
        }
    }

    RevWs = new WebSocket("ws://localhost:5050/chatroom/recive")
    RevWs.onopen=function(){
        console.log("Recive connected");
    }
    RevWs.onmessage = function(e){
        console.log("get message: "+e.data)
    }

    $('#submit').click(function (){
        let message={
            "sendtime":getcurTime(),
            "content":$('#text').val()
        }
        console.log(message)
        $.ajax({
            type: 'post',
            dataType:'json',
            data:message,
            url: '/chatroom/send',
                success: function (res) {
                $('#send').append("<p>"+message.sendtime+" "+message.content+"</p>")
                console.log(res.msg)
            }
        })
    })
    $('#offline').click(function (){
        $.ajax({
            type: 'get',
            url: '/chatroom/offline',
            success: function (res) {
                console.log(res.msg)
                listWs.close(1000,"closed")
                window.location.replace("/")
            }
        })
    })

    $('#erase').click(function (){
        console.log("empty")
        $("#text").val("");
    })
})