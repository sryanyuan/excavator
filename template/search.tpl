<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <title>ExcavatorProxy</title>
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
  <!--search-->
  <nav id="navbar" class="navbar navbar-default navbar-fixed-top" role="navigation">
	<div class="navbar-header">
          <!--button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button-->
          <a id="id-logotext" class="navbar-brand" href="/">ExcavatorProxy</a>
    </div>
			<form id="id-searchform" class="navbar-form form-search active" action="/" method="GET">
				<div class="form-group">
					<input id="id-formsubmit" class="form-control" name="keyword" type="text" value="{{.LastSearch}}" placeholder="搜索" />
					<button id="id-searchbutton" type="submit" class="btn btn-success">搜索</button>
					<i class="fa btn-search fa-search"></i>
					<input id="id-totalpage" type="hidden" name="totalpage" value="{{.TotalPage}}">
					<input id="id-currentpage" type="hidden" name="currentpage" value="{{.CurrentPage}}">
					<input id="id-keyhash" type="hidden" name="keyhash" value="{{.Keyhash}}">
				</div>
				<div id="id-alertgroup" class="alert alert-warning fade in">
					<button id="id-alerttip" type="button" class="close" aria-hidden="true">×</button>  
					<p>Message</p>  
				</div>  
			</form>
	</nav>
  <div class="container">
	<div class="row">
		<div class="col-md-10 col-md-offset-1">
			<!--results-->
			{{if ne .ResultCount 0}}
				<div class="panel panel-default">
					<div class="panel-heading">
						<span>找到{{.ResultCount}}条结果</span>
					</div>
					<div id="id-searchresult" class="panel-body">
						{{if ne .Error ""}}
						<div class="alert alert-danger" role="alert">
							{{.Error}}
						</div>
						{{end}}
						{{range .SearchResult}}
						<div class="alert alert-success" role="alert">
							<strong style="margin-bottom:2px;margin-top:2px;">{{.Name}}</strong>
							<p>创建时间:{{.CollectTime}}  热度:{{.Popular}}  下载速度:{{.DownloadSpeed}}  文件大小:{{.FileSize}}  文件个数:{{.FileCount}}</p>
							<p><a href="{{writeMagnet .}}"  style="word-break:break-all;">磁力链</a></p>
						</div>
						<br/>
						{{end}}
					</div>
					<div id="id-resultcontainerfoot" class="panel-footer hidden">
						<span id="id-scrolltip">下拉至最底部获取下一页数据</span>
						<p><i class="fa fa-angle-double-down"></i></p>
					</div>
				</div>
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
	<p>Just for learning purpose.DO NOT used for commercial purposes <i class="fa fa-copyright"></i> ExcavatorProxy {{.ProcessTime}} ms</p>
  </div>
  <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
  <script src="/static/js/jquery.min.js"></script> 
  <!-- Include all compiled plugins (below), or include individual files as needed -->
  <script src="/static/js/bootstrap.min.js"></script>
  <!-- Custom js -->
  <script src="/static/js/base.js"></script>
  <!-- Adjust footer -->
  <script>
	(function(){
      if($(window).height()==$(document).height()){
        $("#id-footer").addClass("navbar-fixed-bottom");
      }
      else{
        $("#id-footer").removeClass("navbar-fixed-bottom");
      }
    })()
  </script>
  </body>
</html>
