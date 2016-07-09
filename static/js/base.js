//	dynamic load search result
var tplResultItem = '<div class="alert alert-success" role="alert">\
					<strong style="margin-bottom:2px;margin-top:2px;">{{Name}}</strong>\
					<p>创建时间:{{CollectTime}}  热度:{{Popular}}  下载速度:{{DownloadSpeed}}  文件大小:{{FileSize}}  文件个数:{{FileCount}}</p>\
					<p><a href="{{MagnetLink}}"  style="word-break:break-all;">磁力链</a></p>\
				</div>\
				<br/>';
				
$(window).scroll(function () {
    if ($(document).scrollTop() + $(window).height() >= $(document).height()) {
        requestNewPage();
    }
});

var pullingData = false
function requestNewPage() {
	if(pullingData){
			return;
		}
		
		//	pulling data
		var keyhash = $("#id-keyhash").attr("value");
		if(keyhash.length == 0) {
			return;
		}
		var totalPage = Number($("#id-totalpage").attr("value"));
		var currentPage = Number($("#id-currentpage").attr("value")); 
		var nextPage = currentPage + 1;
		
		if(currentPage >= totalPage){
			return;
		}
		
		pullingData = true;
		hideAlertMessage();
		$.get("/page?keyhash="+keyhash+"&page="+nextPage, function(ret){
			//	close the alertgroup
			//hideAlertMessage();
			
			pullingData = false;
			
			var obj = jQuery.parseJSON(ret);
			var resultSize = obj.length;
			
			if(0 == resultSize){
				showAlertMessage("获取数据超时，请重试");
				return;
			}
			
			$("#id-currentpage").attr("value", ""+nextPage);
			
			for(var i = 0; i < resultSize; i++){
				var itemHtml = tplResultItem;
				var dataItem = obj[i];
				itemHtml = itemHtml.replace("{{Name}}", dataItem.Name).replace("{{CollectTime}}", dataItem.CollectTime).replace("{{Popular}}", ""+dataItem.Popular).replace("{{DownloadSpeed}}", dataItem.DownloadSpeed).replace("{{FileSize}}", dataItem.FileSize).replace("{{FileCount}}", dataItem.FileCount).replace("{{MagnetLink}}", dataItem.MagnetURI)
				$("#id-searchresult").append(itemHtml);
			}
			
			if (nextPage >= totalPage) {
				//	hide footer
				$("#id-resultcontainerfoot").find("p").addClass("hidden");
				$("#id-resultcontainerfoot").find("span#id-scrolltip").html("已到达最后一页")
			}
		}).error(function(){
			pullingData = false;
			showAlertMessage("拉取数据失败，请重试");
		});
}

$("#id-alerttip").click(function(){
	//$("#id-alertgroup").fadeOut();
	hideAlertMessage();
})

function showAlertMessage(msg) {
	$("#id-alertgroup").find("p").html(msg);
	//$("#id-alertgroup").fadeIn();
	//$("#id-alerttip").fadeIn();
	$("#id-alertgroup").removeClass("hidden");
}

function hideAlertMessage() {
	//$("#id-alertgroup").hide();
	if (!$("#id-alertgroup").hasClass("hidden")){
		$("#id-alertgroup").addClass("hidden");
	}
}

/*$("#id-searchbutton").click(function(){
	var keyword = $("#id-formsubmit").val();
	if(keyword.length == 0){
		showAlertMessage("关键词不能为空");
		return;
	}
	
	hideAlertMessage();
	$("#id-formsubmit").submit();
})*/

/*$("#id-formsubmit").submit(function(event){
	event.preventDefault();
	
	var keyword = $("#id-formsubmit").val();
	if(keyword.length == 0){
		showAlertMessage("关键词不能为空");
		return;
	}
	
	hideAlertMessage();
	$("#id-formsubmit").submit();
});*/

$(document).ready(function(){
	//	hide alert group
	//$("#id-alertgroup").hide();
	//$("#id-alertgroup").removeClass("hidden");
	
	//	hide ?
	var totalPage = Number($("#id-totalpage").attr("value"));
	var currentPage = Number($("#id-currentpage").attr("value")); 
	
	if (0 != totalPage) {
		$("#id-resultcontainerfoot").removeClass("hidden");
		
		if (currentPage >= totalPage) {
			$("#id-resultcontainerfoot").find("p").addClass("hidden");
			$("#id-resultcontainerfoot").find("span#id-scrolltip").html("已到达最后一页")
		}
	}
});