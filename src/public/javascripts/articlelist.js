Vue.component('article-list', {
  template: `<ol v-bind:style="">
                    <li v-for="item in articlelists">
                        {{item.title}}--{{item.author}}
                    </li>
                </ol>`,
  data: function() {
    return {
      articlelists: []
    }
  },
  methods: {
    getSimpleData: function() {
      this.articlelists.push({
        title: "sdfsdfs",
        author: "liuyingmei"
      })
      let test = this;
      let url = "http://127.0.0.1:8080/article/abstractlist"
      axios.get(url).then(function(response) {
        console.log((response.data)["data"])

        test.articlelists = (response.data)["data"];
      }).catch(function(err) {
        console.log(err)
      })
    },
    setRaw: function(data) {
      this.articlelists = data;
    },
    init: function() {
      let re = []
      re.push({
        title: "sfsfsdf",
        author: "lixiaomeng"
      });
      re.push({
        title: "sdfsdfs",
        author: "liuyingmei"
      })
      console.log(this.data)
      this.articlelists = re;
    }
  }
});
