Vue.component('article-list', {
  template: `<div style="outer_style"><h3 v-bind:style="caption_style">最近的文章</h3>
              <ul v-bind:style="ul_style">
                <li v-for="item in articlelists">
                  <a v-bind:style="link_style" v-bind:href="prefix_article+item.id"><span v-bind:style="span_title_style">{{item.title}}</span><span v-bind:style="span_author_style">{{item.author}}</span></a>
                </li>
              </ul>
              <div v-bind:style="buttonbox_style">
                <button ref="btn_pre" v-on:click="descPage">上翻页</button>
                <button ref="btn_next" v-on:click="ascPage">下翻页</button>
              </div>
              </div>`,
  data: function () {
    return {
      articlelists: [],
      prefix_article: "http://127.0.0.1:8080/view/article/",
      currentpage:0,
      link_style:{
        "text-decoration":"none",
      },
      caption_style:{
        "text-align":"center"
      },
      ul_style:{
        "list-style-type":"square"
      },
      buttonbox_style:{
        "width":"100%",
        "text-align":"center"
      },
      outer_style:{
        "width":"120px"
      },
      span_title_style:{
        "font-size":"larger"
      },
      span_author_style:{
        "font-style":"italic",
        "margin-left":"10px",
        "font-size":"14px"
      }
    }
  },
  mounted: function () {
    console.log("mounted")
    this.$refs.btn_pre.disabled = true;
    this.getSimpleData();
  },
  watch:{
    currentpage:function(val,oldVal){
      if(val==0){
        this.$refs.btn_pre.disabled=true;
      }
    }
  },
  methods: {
    getSimpleData: function (pageno=0) {
      let test = this;
      let url = `http://127.0.0.1:8080/article/abstractlist?pageno=${pageno}`;
      axios.get(url).then(function (response) {
        console.log((response.data)["data"])
        test.articlelists = (response.data)["data"];
      }).catch(function (err) {
        console.log(err)
      })
    },
    setRaw: function (data) {
      this.articlelists = data;
    },
    
    descPage:function(){

      this.currentpage=this.currentpage-1;
      if(this.currentpage<0){
        this.currentpage=0;
        return
      }
      this.getSimpleData(this.currentpage)
    },
    ascPage:function(){
      this.currentpage=this.currentpage+1;
      if(this.articlelists.length<10){
        this.$refs.btn_next.disabled=true;
        return
      }
      this.getSimpleData(this.currentpage)
    }
  }
});
