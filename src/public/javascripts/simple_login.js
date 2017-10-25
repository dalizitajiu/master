Vue.component('simple-login', {
    template: `<div>
                    <form action="/user/login" method="post" enctype="application/x-www-form-urlencoded">
                    请输入邮箱:<br>
                    <input v-bind:value="email" name="email"></input>
                    <br>
                    请输入密码:<br>
                    <input name="pwd" type="password"></input>
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
        }
    }
})