Vue.component('simple-login', {
    template:`
        <div v-bind:style="outer_style">
            <span ref="mspan" v-bind:style="span_style">账户</span><input ref=minput v-bind:style="input_style" v-model.trim="email" name="email" autofocus></input>
            <br>
            <span v-bind:style="span_style">密码</span><input v-bind:style="input_style" v-model.trim="pwd" name="pwd" type="password"></input>
            <br>
            <input ref="btn" v-bind:style="btn_style" type="button" v-on:click="login" value="登录"></input>
        </div>
    `,
    data: function() {
        return {
            email: "",
            pwd: "",
            input_style:{
                "border-radius": "5px",
                "text-align":"center"
            },
            span_style:{
                "margin-right":"10px",
                "color":"#4CAF50"
            },
            btn_style:{
                "border-radius": "5px",
                "color":"#4CAF50",
                "margin-top":"5px"
            },
            outer_style:{
                "width":"300px",
                "height":"100px"
            }
        }
    },
    updated:function(){
        console.log(this.$refs.mspan.offsetWidth,this.$refs.minput.offsettWidth)
        let width=this.$refs.mspan.offsetWidth+this.$refs.minput.offsetWidth+10;
        console.log(width)
        this.$data.outer_style["width"]=`${width}px`;
        let btn_width=this.$refs.btn.clientWidth;
        console.log(btn_width)
        this.$data.btn_style["margin-left"]=`${(width-btn_width)/2}px`;
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