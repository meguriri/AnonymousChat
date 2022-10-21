//获取当前时间
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

    //创建用户列表的websocket连接
    listWs = new WebSocket("ws://localhost:5050/chatroom/userlist")

    //用户列表的websocket建立连接时
    listWs.onopen=function(){
        console.log("userlist connected");
    }

    //用户列表的websocket收到消息时
    listWs.onmessage = function(e){
        let user=JSON.parse(e.data)
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

    //用户列表的websocket关闭时
    listWs.onclose=function (){
        console.log("userlist disconnected")
    }

    //创建消息的websocket连接
    RevWs = new WebSocket("ws://localhost:5050/chatroom/recive")

    //消息的websocket建立连接时
    RevWs.onopen=function(){
        console.log("Receive connected");
    }

    //消息的websocket收到消息时
    RevWs.onmessage = function(e){
        console.log("get message: "+e.data)
    }

    //消息的websocket关闭时
    RevWs.onclose=function (){
        console.log("Receive disconnected");
    }

    //发送信息
    $('#submit').click(function (){
        //创造message实例
        let message={
            "sendtime":getcurTime(),//发送时间
            "content":$('#text').val()//内容
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

    //下线
    $('#offline').click(function (){
        $.ajax({
            type: 'get',
            url: '/chatroom/offline',
            success: function (res) {
                console.log(res.msg)
                //关闭websocket连接
                listWs.close(1000,"closed")
                //返回登录界面
                window.location.replace("/")
            }
        })
    })

    //清除输入框
    $('#erase').click(function (){
        console.log("empty")
        $("#text").val("");
    })

})