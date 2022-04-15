$(document).ready(function () {
    listWs = new WebSocket("ws://localhost:5050/chatroom/userlist")
    listWs.onopen=function(){
        console.log("connected");
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

    $('#erase').click(function (){
        console.log("empty")
        $("#text").val("");
    })
})