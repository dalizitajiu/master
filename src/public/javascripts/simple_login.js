Vue.component('simple-login', {
    template:`
        <div>
            <span stype="span_style">账户</span><input style="input_style" v-model.trim="email" name="email" autofocus></input>
            <br>
            <span stype="span_style">密码</span><input v-model.trim="pwd" name="pwd" type="password"></input>
            <br>
            <input type="button" v-on:click="login" value="登录"></input>
        </div>
    `,
    data: function() {
        return {
            email: "",
            pwd: "",
            input_style:{
                "border-radius": "5px",
            },
            span_style:{
                "margin-right":"10px",
                "color":"#4CAF50"
            }
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