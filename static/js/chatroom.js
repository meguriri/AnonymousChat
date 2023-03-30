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

function getGender(gender){
    if(gender===0){
        return ' ♂ '
    }
    return ' ♀ '
}

var listWs
var RevWs
$(document).ready(function () {

    //获取聊天记录功能
    $.ajax({
        type: 'get',
        url: '/chatroom/chats',
        success: function (res) {
            console.log(res)
            for(let i=0;i<res.message.length;i++){
                let message=res.message[i]
                console.log(message)
                let gender=getGender(message.senduser.gender)

                $('#chatroom').append('<div class="row p-1 mb-2">\n' +
                    '                        <div class="col-1 m-0 p-0">\n' +
                    '                            <li class="my-1 border-0 list-group-item d-flex">\n' +
                    '                                <span class=" rounded-circle text-white text-center" style="width: 50px;height: 50px;\n' +
                    '                                background-color: rgb('+message.senduser.color[0]+','+message.senduser.color[1]+','+message.senduser.color[2]+')"><i class="bi bi-person" style="font-size: 32px;"></i> </span>\n' +
                    '                            </li>\n' +
                    '                        </div>\n' +
                    '                        <div class="col-4 m-0 p-0">\n' +
                    '                            <h5 class="text-black-50">'+message.senduser.nickname+gender+'   '+message.sendtime+'</h5>\n' +
                    '                            <div class="p-2 border" style="background-color: rgb(233, 236, 236);display:inline-block; word-break: break-all;word-wrap: break-word;border-radius: 10px;">\n' +
                    '                                <p style="font-size:20px;letter-spacing: 1px;"><b>'+message.content +
                    '                                </b></p>\n' +
                    '                            </div>\n' +
                    '                        </div>\n' +
                    '                    </div>')
            }
        }
    })

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
            c=getGender(user[i].gender)
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
    RevWs = new WebSocket("ws://localhost:5050/chatroom/recive")

    //消息的websocket建立连接时
    RevWs.onopen=function(){
        window.alert("欢迎来到匿名聊天室，注意文明用语。（注意：超过10分钟未发送消息将会自动踢出聊天室！！！）")
        console.log("Receive connected");
    }

    //消息的websocket收到消息时
    RevWs.onmessage = function(e){
        console.log("get message: "+e.data)
        let message=JSON.parse(e.data)
        let gender=getGender(message.senduser.gender)
        $('#chatroom').append('<div class="row p-1 mb-2">\n' +
            '                        <div class="col-1 m-0 p-0">\n' +
            '                            <li class="my-1 border-0 list-group-item d-flex">\n' +
            '                                <span class=" rounded-circle text-white text-center" style="width: 50px;height: 50px;\n' +
            '                                background-color: rgb('+message.senduser.color[0]+','+message.senduser.color[1]+','+message.senduser.color[2]+')"><i class="bi bi-person" style="font-size: 32px;"></i> </span>\n' +
            '                            </li>\n' +
            '                        </div>\n' +
            '                        <div class="col-4 m-0 p-0">\n' +
            '                            <h5 class="text-black-50" >'+message.senduser.nickname+gender+'   '+message.sendtime+
            '                            </h5><div class="p-2 border" style="background-color: rgb(233, 236, 236);display:inline-block; word-break: break-all;word-wrap: break-word;border-radius: 10px;">\n' +
            '                                <p style="font-size:20px;letter-spacing: 1px;"><b>'+message.content+'</b></p>\n' +
            '                            </div>\n' +
            '                        </div>\n' +
            '                    </div>')
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
        let text=$('#text').val().toString()
        console.log(text)
        if(text.indexOf("<")!==-1||text.indexOf("</")!==-1||text.indexOf(">")!==-1){
            alert("不可以插入脚本！")
            return
        }
        let message={
            "sendtime":getcurTime(),//发送时间
            "content":text//内容
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
                let gender=''
                if(user.gender==0){
                    gender=' ♂ '
                }else{
                    gender=' ♀ '
                }
                $('#chatroom').append('<div class="row p-1 mb-2" >\n' +
                    '                        <div class="col-4  offset-7">\n' +
                    '                            <h5 class="text-black-50 " style="text-align:right;">'+user.nickname+gender+'   '+message.sendtime+
                    '                            </h5><div class="row"><div class="col-6 p-2 ms-auto" style="text-align:right;background-color: rgb(230, 245, 255);display:inline-block; word-break: break-all;word-wrap: break-word;border-radius: 10px;">\n' +
                    '                                <p style="font-size:20px;letter-spacing: 1px;"><b>'+message.content+'</b></p>\n' +
                    '                            </div></div>\n' +
                    '                        </div>\n' +
                    '                        <div class="col-1 m-0 p-0">\n' +
                    '                            <li class="my-1 border-0 list-group-item d-flex">\n' +
                    '                                <span class=" rounded-circle text-white text-center" style="width: 50px;height: 50px;\n' +
                    '                                background-color: rgb('+user.color[0]+','+user.color[1]+','+user.color[2]+')"><i class="bi bi-person" style="font-size: 32px;"></i> </span>\n' +
                    '                            </li>\n' +
                    '                        </div>\n' +
                    '                    </div>')
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
                $.removeCookie('login', { path: '/' })
                window.location.replace("/")
            }
        })
    })

    //清除输入框
    $('#erase').click(function (){
        console.log("empty")
        $("#text").val("")
    })

})