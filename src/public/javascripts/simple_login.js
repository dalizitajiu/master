Vue.component('simple-login', {
    // template: `<div>
    //                 <form ref="loginform" action="/user/login" method="post" enctype="application/x-www-form-urlencoded" v-on:submit.prevent="checkForm">
    //                 请输入邮箱:<br>
    //                 <input v-model.trim="email" name="email"></input>
    //                 <br>
    //                 请输入密码:<br>
    //                 <input v-model.trim="pwd" name="pwd" type="password"></input>
    //                 <br>
    //                 <input type="submit"></input>
    //                 </form>
    //             </div>`,
    template:`
        <div>
            <input v-model.trim="email" name="email"></input>
            <br>
            <input v-model.trim="pwd" name="pwd" type="password"></input>
            <br>
            <input type="button" v-on:click="login" value="登录"></input>
        </div>
    `,
    data: function() {
        return {
            email: "",
            pwd: ""
        }
    },
    methods: {
        setEmail: function(newemail) {
            this.email = newemail.toString();
        },
        getEmail: function() {
            return this.email;
        },
        setPasswd: function(newpwd) {
            this.pwd = newpwd;
        },
        getPasswd: function() {
            return this.pwd;
        },
        checkForm:function(){
            console.log("dagoushi")
            if(this.email.length<1 || this.pwd.length<1){
                console.log("shit")
                return false;
            }else{
                this.$refs.loginform.submit();
                return true;
            }
        },
        login:function () {
            let param = new URLSearchParams()
            param.append("email",this.email)
            param.append("pwd",this.pwd)
            axios.post("http://127.0.0.1:8080/user/login",param).then(function(response){
                let data=response.data;
                if(data.errno==0){
                    window.location.href="http://127.0.0.1:8080/view/article/getones"
                }else{
                    alert("登录失败")
                }
                
            }).catch(function(err){
                console.log(err)
            })
        }

    }
})