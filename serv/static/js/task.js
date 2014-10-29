var task = (function($){
    //初始化
    function init(containID, url){
	taskData(containID, url);
	$("#newTypeTxt").focus(function(){ 
	    $("#newTypeTxt").val("").css("border-color", "#ccc");
	});
	$("#searchTxt").focus(function(){ 
	    $("#searchTxt").val("").css("border-color", "#ccc");
	});

    }

    //任务数据请求
    function taskData(containID, url){
	$.ajax({
	    type: 'GET',
	    url: url,
	    async: true,
	    success: function(data){
		var json, datatable = "", ttype = "";
		if (data != null || data != "") {
		    json = JSON.parse(data);
		    for (var i = 0, len = json.length; i < len; i++) {
			//加载数据表
			datatable += docDataModule(json[i]);
			ttype += '<option value="' + json[i].Name + '">' + json[i].Name + '</option>';
		    }
		    $("#" + containID).html(datatable);
		    $("#select").html(ttype);
		    $("#search").click(function(){searchInfo('/monitor?act=info'); });
		    $("#newType").click(function(){ addType('/newtype'); });


		    alertInfo('/monitor?act=info');
		}
	    },
	    error: function(data) {
		$("#" + containID).html('ERR...' + data.responseText);
	    }
	});
    }

    function redodel(url, ttype, tid){
	$.ajax({
		type: 'GET',
		url: url+ '&ttype=' + ttype + '&tid=' + tid,
		success: function(data){
		    alert(data);
		    location.reload();
		},
		error: function(msg) {
		    alert(msg.responseText + ". stat:" + msg.status)
		}
	});
    }

    function searchInfo(url){
	var index = $("#select")[0].selectedIndex;
	var ttype = $("#select")[0].options[index].text;
	var tid = $("#searchTxt").val();
	if (tid != "") {
	    $.ajax({
		type: 'GET',
		url: url + '&ttype=' + ttype + '&tid=' + tid,
		success: function(data){
		    var json;
		    if (data != "") {
			console.info(data == "");
			json = JSON.parse(data);
			Alert('<table class="info">' + docInfo(json) + '</table>');
		    } else if (data == "") {
			Alert('未找到该任务信息');
		    }
		},
		error: function(msg) {
		    alert(msg.responseText + ". stat:" + msg.status)
		}
	    });
	} else {
	    $("#searchTxt").css("border-color", "#f30");
	}
    }

    function addType(url){
	var name = $("#newTypeTxt").val();
	if (/^[a-z0-9]{3,6}$/.test(name)) {
	    if (confirm("确定添加数据")) {
		$.ajax({
		    type: 'GET',
		    url: url + '?name=' + name,
		    success: function(data){
		    	location.reload();
		    },
		    error: function(msg) {
			alert(msg.responseText + ". stat:" + msg.status)
		    }
		});
	    }
	} else {
	    $("#newTypeTxt").css("border-color", "#f30");

	}
    }

    //doc aside
    function insertIntoAside(name, ibuff, obuff, ncout, icout, scout, ecout){
	var string = '<aside class="aside box-flex1">' +
	    '<p class="task-type">任务类型 : ' + name + '</p>' +
	    '<p class="asiderow"><label>入缓存 </label>: ' + ibuff + '</p>' +
	    '<p class="asiderow"><label>出缓存 </label>: ' + obuff + '</p>' +
	    '<p class="asiderow"><label>新增 </label>: ' + ncout + '</p>' +
	    '<p class="asiderow"><label>进行 </label>: ' + icout + '</p>' +
	    '<p class="asiderow"><label>成功 </label>: ' + scout + '</p>' +
	    '<p class="asiderow"><label>失败 </label>: ' + ecout + '</p>' +
	    '</aside>';

	return string;
    }

    //doc
    function insertIntoDataTD(arr, name, type){
	var string = '';
	if (arr instanceof Array) {
	    for (var i = 0, len = arr.length; i < len; i++) {
	    	if (type == 'rapper') {
	    		string += '<li ttype="' + name + '" class="rapper">' + arr[i] + '</li>';
	    		continue;
	    	}
	    	if (type == 're') {
	    		string += '<li ttype="' + name + '" type="redo">' + arr[i] + '</li>';
	    		continue;
	    	}
		string += '<li ttype="' + name + '">' + arr[i] + '</li>';
	    }
	}

	return string;
    }
    //doc module
    function docDataModule(options){
	var string = '<div class="dataModule box">' +
	    insertIntoAside(options.Name, options.Ibuff, options.Obuff, options.Ncout, options.Icout, options.Scout, options.Ecout) +
	    '<div class="datatable box">' +
	    '<ul class="new box-flex1">' +
	    insertIntoDataTD(options.Nrec, options.Name) +
	    '</ul>' +
	    //ing
	    '<ul class="ing box-flex1">' +
	    insertIntoDataTD(options.Irec, options.Name) +
	    '</ul>' +
	    //error
	    '<ul class="error box-flex1">' +
	    insertIntoDataTD(options.Erec, options.Name, 're') +
	    '</ul>' +
	    //rapper
	    '<ul class="rapper box-flex1">' +
	    insertIntoDataTD(options.Rappers, options.Name, 'rapper') +
	    '</ul>' +
	    '</div>' +
	    '</div>';

	return string;
    }

    //click detail info
    function alertInfo(url){
	$(".datatable li").not('.rapper').click(function(){
	    var ttype = $(this).attr('ttype');
	    var tid = $(this).html();
	    var _this = $(this);
	    $.ajax({
		type: 'GET',
		url: url + '&ttype=' + ttype + '&tid=' + tid,
		success: function(data){
		    var json;
		    if (data != null || data != "") {
			json = JSON.parse(data);
			Alert('<table class="info">' + docInfo(json) + '</table>');
			var type = _this.attr('type');
			if (type == 'redo') {
				$("#redo_btn").show();
				$("#del_btn").show();
			} else {
				$("#redo_btn").hide();
				$("#del_btn").hide();
			}

			$("#redo_btn").click(function(){redodel('/monitor?act=redo', ttype, tid); });
			$("#del_btn").click(function(){ redodel('/monitor?act=del', ttype, tid); });
		    }
		},
		error: function(data){
		    
		}
	    });
	});
    }

    //info frame
    function docInfo(data){
	var string = '';
	var obj = {
	    tid: '任务序号',
	    rid: '记录序号',
	    info: '任务内容',
	    stat: '任务状态',
	    addTime: '添加时间',
	    overTime: '完成时间',
	    rapper: '工兵',
	    client: '客户端',
	    remark: '备注'
	};
	for (var prop in obj) {
	    string += '<tr class="clearfix">' +
		'<td class="infoTit">' + obj[prop] + '</td>' + 
		'<td class="infoCon">' + data[prop] + '</td>' +
		'</tr>';
	}

	return string;
    }

    //弹框组件
    function Alert(txt) {
	var dg = function(id){ return document.getElementById(id); },  //封装ID
	txt = txt,  //接受文本
	forms = dg("modal"), //弹出框
	layer = document.createElement('div'), //创建悬浮层dom
	close = [dg("modal_close"), dg("modal_btn"), layer]; //关闭按钮DOM数组
	//为悬浮层定义calss
	layer.className = "fade";

	//事务操作
	//弹出信息
	dg("modal_txt").innerHTML = txt;
	//弹出框显示
	forms.className = "modal fade in";
	//文档中插入悬浮层
	document.body.appendChild(layer);
	//悬浮层淡入
	layer.className = "modal_backdrop fade in";

	//关闭数组中dom点击淡出悬浮层和弹出框
	for (var i = 0; i < close.length; i++) {
	    close[i].onclick = function(){
		layer.className = "modal_backdrop fade";
		forms.className = "modal fade";
		setTimeout(function(){
		    document.body.removeChild(layer);
		}, 80);
	    }
	}
    }

    return {
	init: init
    }
})(Zepto);

function Alert(txt) {
	var dg = function(id){ return document.getElementById(id); },  //封装ID
		txt = txt,  //接受文本
		forms = dg("modal"), //弹出框
		layer = document.createElement('div'), //创建悬浮层dom
		close = [dg("modal_close"), dg("modal_btn"), layer]; //关闭按钮DOM数组
	//为悬浮层定义calss
	layer.className = "fade";

	//事务操作
	//弹出信息
	dg("modal_txt").innerHTML = txt;
	//弹出框显示
	forms.className = "modal fade in";
	//文档中插入悬浮层
	document.body.appendChild(layer);
	//悬浮层淡入
	layer.className = "modal_backdrop fade in";

	//关闭数组中dom点击淡出悬浮层和弹出框
	for (var i = 0; i < close.length; i++) {
		close[i].onclick = function(){
			layer.className = "modal_backdrop fade";
			forms.className = "modal fade";
			setTimeout(function(){
				document.body.removeChild(layer);
			}, 200);
		}
	}
}

task.init('insert', '/monitor');
