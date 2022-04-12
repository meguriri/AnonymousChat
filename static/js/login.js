$(document).ready(function () {
    $('#submitcolor').click(function(){
        let c1=$('#c1').val()
        let c2=$('#c2').val()
        let c3=$('#c3').val()
        $('#123').css("background-color", "rgb("+String(c1)+","+String(c2)+","+String(c3)+")")
    })
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
    $('#submit').click(function (){
        console.log("click..")
        let user ={
            nickname: $('#nicknameinput').val(),
            gender: Number($("input[type='radio']:checked").val()),
            color: [
                parseInt($('#c1').val()),
                parseInt($('#c2').val()),
                parseInt($('#c3').val())
            ]
        }
        let json=JSON.stringify(user)
        $.ajax({
            type: 'post',
            url: '/login',
            dataType: 'json',
            data: json,
            success: function (res) {
                console.log(res.msg)
                window.location.replace("/chatroom/")
            }
        })
    })
})