Vue.component('my-articles', {
    template: `<div><h3 v-bind:style="caption_style">我的全部文章</h3>
              <ol v-bind:style="">
                <li v-for="item in articlelists">
                  <a v-bind:style="link_style" v-bind:href="prefix_article+item.id">{{item.title}}--{{item.author}}</a>
                </li>
              </ol>
              </div>`,
    data: function () {
        return {
            articlelists: [],
            prefix_article: "http://127.0.0.1:8080/view/article/",
            link_style: {
                "text-decoration": "none",
            },
            caption_style: {
                "text-align": "center"
            },
            buttonbox_style: {
                "width": "100%",
                "text-align": "center"
            }
        }
    },
    methods: {
        getSimpleData: function (pageno = 0) {
            let test = this;
            let url = `http://127.0.0.1:8080/article/getones`;
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
        init: function () {
            let re = [];
            re.push({
                title: "sfsfsdf",
                author: "lixiaomeng",
                link: "https://www.baidu.com"
            });
            re.push({
                title: "sdfsdfs",
                author: "liuyingmei",
                link: "http://www.baidu.com"
            });
            console.log(this.data)
            this.articlelists = re;
        }
    }
});
