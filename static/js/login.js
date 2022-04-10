$(document).ready(function () {
    $('#enter').click(function (){
        console.log("click..")
        let a =$('#user').val()
        $.ajax({
            type: 'post',
            url: '/login',
            dataType: 'json',
            data: {name: a},//序列化表单值
            success: function (res) {
                console.log(res.msg)
            }
        })
    })
})