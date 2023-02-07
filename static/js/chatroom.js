//获取当前时间
function getcurTime() {
    let Digital=new Date();
    let hours=Digital.getHours();
    let minutes=Digital.getMinutes();
    let seconds=Digital.getSeconds();
    if(minutes<=9){
        minutes="0"+minutes;
    }if(seconds<=9){
        seconds="0"+seconds;
    }
    myclock=hours+":"+minutes+":"+seconds;
    return myclock;
}
var listWs
var RevWs
$(document).ready(function () {
    //创建用户列表的websocket连接
    listWs = new WebSocket("ws://192.168.31.177:5050/chatroom/userlist")

    //用户列表的websocket建立连接时
    listWs.onopen=function(){
        console.log("userlist connected");
        alert("userlist connected")
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
        // $.removeCookie('login', { path: '/' })
        // window.location.replace("/")
    }

    //创建消息的websocket连接
    RevWs = new WebSocket("ws://192.168.31.177:5050/chatroom/recive")

    //消息的websocket建立连接时
    RevWs.onopen=function(){
        window.alert("欢迎来到匿名聊天室，注意文明用语。（注意：超过10分钟未发送消息将会自动踢出聊天室！！！）")
        console.log("Receive connected");
        alert("Receive connected")
    }

    //消息的websocket收到消息时
    RevWs.onmessage = function(e){
        console.log("get message: "+e.data)
        let message=JSON.parse(e.data)
        $('#receive').append("<p>"+message.senduser.nickname+": "+message.sendtime+" "+message.content+"</p>")
    }

    //消息的websocket关闭时
    RevWs.onclose=function (){
        console.log("Receive disconnected")
        $.removeCookie('login', { path: '/' })
        alert("由于您长时间未发送消息，系统自动下线。请重新登录！")
        window.location.replace("/")
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
                console.log(res)
                let user =JSON.parse(res.user)
                $('#send').append("<p>"+user.nickname+": "+message.sendtime+" "+message.content+"</p>")
                $("#text").val("")
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
            }
        })
    })

    //清除输入框
    $('#erase').click(function (){
        console.log("empty")
        $("#text").val("")
    })

})