const axios = require("axios");
const qs=require("querystring");
const login = async function(username, pwd) {
	const url = "http://api.u.panda.tv:8360/test/login";
	let param={
		"username":username,
		"pwd":pwd
	}
	const response = await axios.post(url, qs.stringify(param));
	console.log(response.data);
	return response.data;
}
// $.ajax({
// 		url: '/path/to/file',
// 		type: "POST"
// 		dataType: "json"
// 		data: {
// 			anthor: 'lixiaomeng',
// 			title:"dsf",
// 			subtitle:"",
// 			content:tinyMCE.activeEditor.getContent(),
// 			ttime:"",
// 			ttoken:""
// 		},
// 	})
// 	.done(function() {
// 		console.log("success");
// 	})
// 	.fail(function() {
// 		console.log("error");
// 	})
// 	.always(function() {
// 		console.log("complete");
// 	});
//