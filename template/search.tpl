<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <title>Excavator search magnet link</title>
    <!-- Bootstrap -->
    <link href="/static/css/bootstrap.min.css" rel="stylesheet" />
    <!-- Bootstrap theme -->
    <link href="/static/css/bootstrap-theme.min.css" rel="stylesheet" />
	<!-- Font awesome -->
	<link href="/static/css/font-awesome.min.css" rel="stylesheet" />
	<!-- Custom css -->
    <link href="/static/css/base.css" rel="stylesheet" />
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="//cdn.bootcss.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="//cdn.bootcss.com/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
  <div class="container">
	<div class="row">
		<div class="col-md-10 col-md-offset-1">
			<form class="navbar-form form-search active" action="/" method="GET">
				<div class="form-group">
					<input class="form-control" name="keyword" type="text" value="{{.LastSearch}}" placeholder="搜索" />
				</div>
				<i class="fa btn-search fa-search"></i>
				<span>抓取页数：</span>
				<select name="maxpage">
					<option value ="1">1</option>
					<option value ="5">5</option>
					<option value="10">10</option>
					<option value="20">20</option>
					<option value="40">40</option>
					<option value="0">all</option>
				</select>
			</form>
			
			{{if ne .ResultCount 0}}
				<span style="margin-bottom:4px;">找到{{.ResultCount}}条结果</span>
				{{range .SearchResult}}
				<div class="alert alert-success" role="alert">
					<strong style="margin-bottom:2px;margin-top:2px;">{{.Name}}</strong>
					<p>创建时间:{{.CollectTime}}  热度:{{.Popular}}  下载次数:{{.DownloadTimes}}  下载速度:{{.DownloadSpeed}}  文件大小:{{.FileSize}}  文件个数:{{.FileCount}}</p>
					<p><a href="{{writeMagnet .}}"  style="word-break:break-all;">磁力链</a></p>
				</div>
				<br/>
				{{end}}
			{{else}}
				{{if ne .LastSearch ""}}
			
				{{if ne .Error ""}}
				<div class="alert alert-danger" role="alert">
					{{.Error}}
				</div>
				{{else}}
				<div class="alert alert-success" role="alert">
					没有找到结果
				</div>
				{{end}}
				{{end}}
			{{end}}
		</div>
	</div>
  </div>
  
  <div id="id-footer">
	<p>Dev by sryan. Powered by golang and bootstrap</p>
	<p>Just for learning purpose.DO NOT used for commercial purposes </p>
  </div>
  <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
  <script src="/static/js/jquery.min.js"></script> 
  <!-- Include all compiled plugins (below), or include individual files as needed -->
  <script src="/static/js/bootstrap.min.js"></script>
  <!-- Custom js -->
  <script src="/static/js/base.js"></script>
  <!-- Adjust footer -->
  <!--script>
	(function(){
      if($(window).height()==$(document).height()){
        $("#id_footer").addClass("navbar-fixed-bottom");
      }
      else{
        $("#id_footer").removeClass("navbar-fixed-bottom");
      }
    })()
  </script-->
  </body>
</html>
