Vue.component("article", {
    template: `<div v-bind:style="outer_style">
                    <h1 v-bind:style="title_style">{{title}}</h1>
                    <h3 v-bind:style="author_style">{{author}}</h3>
                    <h5 v-bind:style="title_style">{{createtime}}</h5>
                    <div ref="content" v-bind:style="content_style"></div>
                </div>`,
    data: function () {
        return {
            author: "jack",
            title: "are u ok",
            content: "no thing to be done",
            createtime: "2017/10/27",
            id:17,
            outer_style:{
                "width":"60%",
                "margin-left":"20%",
            },
            title_style:{
                "text-align":"center",
            },
            author_style:{
                "text-align":"left",
                "margin-left":"70%",
            },
            content_style:{
                "font-size":"14px"
            }
        }
    },
    created:function(){
        console.log("created");
        const myreg=RegExp(/\/view\/article\/(\d+)$/);
        let path=window.location.pathname;
        let reslist = myreg.exec(path);
        this.id=reslist[1];
        console.log(this.id);
    },
    mounted:function(){
        this.getContent();
    },
    methods: {
        getContent: function () {
            let self=this;
            let url=`http://127.0.0.1:8080/article/${this.id}`;
            axios.get(url).then(function(response){
                if(response.data.errno!=0){
                    return
                }
                let res=response["data"]["data"]
                self.author=res.author;
                self.title=res.title;
                self.content = this.unescapeString(res.content);
                self.createtime=res.createtime;
                this.$refs.content.innerHTML=self.content;
            })
        },
        unescapeString:function(raw){
            return raw.replace(/(\&|\&)gt;/g, ">").replace(/(\&|\&)lt;/g, "<").replace(/(\&|\&)quot;/g, "\"");
        }

    }
})