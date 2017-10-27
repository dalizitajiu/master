$.ajax({
		url: '/path/to/file',
		type: "POST"
		dataType: "json"
		data: {
			anthor: 'lixiaomeng',
			title:"dsf",
			subtitle:"",
			content:tinyMCE.activeEditor.getContent(),
			ttime:"",
			ttoken:""
		},
	})
	.done(function() {
		console.log("success");
	})
	.fail(function() {
		console.log("error");
	})
	.always(function() {
		console.log("complete");
	});