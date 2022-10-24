$(document).ready(function () {

    //用户头像颜色
    $('#submitcolor').click(function(){
        let c1=$('#c1').val()
        let c2=$('#c2').val()
        let c3=$('#c3').val()
        $('#123').css("background-color", "rgb("+String(c1)+","+String(c2)+","+String(c3)+")")
    })

    //获取在线人数
    $('#onlineusers').each(function (){
        $.ajax({
            type: 'get',
            url: '/onlineusers',
            success: function (res) {
                console.log(res.onlineusers)
                $('#onlineusers').text(String(res.onlineusers))
            }
        })
    })

    //用户登录
    $('#submit').click(function (){
        console.log("click..")
        //用户对象
        let user ={
            nickname: $('#nicknameinput').val(),
            gender: Number($("input[type='radio']:checked").val()),
            color: [
                parseInt($('#c1').val()),
                parseInt($('#c2').val()),
                parseInt($('#c3').val())
            ]
        }
        if(user.nickname.length>10||user.nickname.length<1){
            $('#nicknameinput').val("")
            window.alert("输入昵称不规范，请重新输入昵称！")
        }else{
            let json=JSON.stringify(user)
            $.ajax({
                type: 'post',
                url: '/login',
                dataType: 'json',
                data: json,
                success: function (res) {
                    if(res.msg==="ok"){
                        window.location.replace("/chatroom/")
                    }else{
                        $('#nicknameinput').val("")
                        window.alert("该名称已被使用，请重新输入昵称！")
                    }
                }
            })
        }
    })

})