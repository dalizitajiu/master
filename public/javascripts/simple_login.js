Vue.component('simple-login', {
    template: `<div><input v-bind:value="email"></input><input v-bind:value="pwd" type="password"></input></div>`,
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