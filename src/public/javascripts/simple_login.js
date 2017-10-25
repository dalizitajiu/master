Vue.component('simple-login', {
    template: `<div>
                    <form ref="loginform" action="/user/login" method="post" enctype="application/x-www-form-urlencoded" v-on:submit.prevent="checkForm">
                    请输入邮箱:<br>
                    <input v-model.trim="email" name="email"></input>
                    <br>
                    请输入密码:<br>
                    <input v-model.trim="pwd" name="pwd" type="password"></input>
                    <br>
                    <input type="submit"></input>
                    </form>
                </div>`,
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
        }
    }
})