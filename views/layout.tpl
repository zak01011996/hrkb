<!DOCTYPE html>

<html>
<head>
  <title>{{.title}}</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
  <!--CSS-->
  <link rel="stylesheet" type="text/css" href="/static/bower_components/jquery-ui/themes/redmond/jquery-ui.min.css">
  <link rel="stylesheet" type="text/css" href="/static/bower_components/bootstrap/dist/css/bootstrap.min.css">
  <link type="text/css"href="/static/bower_components/bootstrap-material-design/dist/css/ripples.min.css" rel="stylesheet">
  <link type="text/css" href="/static/bower_components/bootstrap-material-design/dist/css/material-wfont.min.css" rel="stylesheet">
  <link rel="stylesheet" type="text/css"  href="/static/bower_components/snackbarjs/dist/snackbar.min.css">
  <link rel="stylesheet" type="text/css"  href="/static/bower_components/dropzone/dist/min/basic.min.css">
  <link rel="stylesheet" type="text/css"  href="/static/bower_components/dropzone/dist/min/dropzone.min.css">
  <link rel="stylesheet" type="text/css" href="/static/css/main.css">
  {{range $style := .styles}}
	{{$style | css}}
  {{end}}
  <!--/CSS-->
  <!--JS-->
  <script src="/static/bower_components/jquery/dist/jquery.min.js"></script>
  <script src="/static/bower_components/bootstrap/dist/js/bootstrap.min.js"></script>
  <script src="/static/bower_components/bootstrap-material-design/dist/js/material.min.js"></script>
  <script src="/static/bower_components/bootstrap-material-design/dist/js/ripples.min.js"></script>
  <script src="/static/bower_components/snackbarjs/dist/snackbar.min.js"></script>
  <script src="/static/bower_components/dropzone/dist/min/dropzone.min.js"></script>
  <script src="/static/bower_components/jquery-ui/ui/core.js"></script>
  <script src="/static/bower_components/jquery-ui/ui/widget.js"></script>
  <script src="/static/bower_components/jquery-ui/ui/menu.js"></script>
  <script src="/static/bower_components/jquery-ui/ui/position.js"></script>
  <script src="/static/bower_components/jquery-ui/ui/autocomplete.js"></script>
  <script src="/static/bower_components/jquery.cookie/jquery.cookie.js"></script>
  <script src="/static/bower_components/lodash/lodash.min.js"></script>
  <script>
  var t_save='{{T "save"}}';
  var t_cancel='{{T "cancel"}}';
  </script>
  <script src="/static/js/hrkb.js"></script>
  <script src="/static/js/ratings.js"></script>
  <script>
      $(function() {
          $.material.init();
      })
  </script>
	
	{{range $js := .scripts }}
		{{$js | js}}  
  {{end}}
    <meta name="theme-color" content="#009587">
 <!--/JS-->
</head>

<body>
<!-- menu -->
{{ if .UserName }}
<div class="navbar navbar-default">
    <div class="navbar-header">
        <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-responsive-collapse">
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
        </button>
	{{ $url := urlFor "Cand.Index" }}
        <a class="navbar-brand" href="{{ $url }}">HRM</a>
    </div>
    <div class="navbar-collapse collapse navbar-responsive-collapse">
        <ul class="nav navbar-nav">
	    {{if .UserIsAdmin}}
	    {{ $url := urlFor "User.Index" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{ T "user" 2}}</a></li>
	    {{end}}

	    {{ $url := urlFor "Dep.Index" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{ T "dep" 2}}</a></li>
	    {{ $url := urlFor "Crit.Index" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{ T "crit" 2}}</a></li>
	    {{ $url := urlFor "Cand.Index" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{ T "cand" 2}}</a></li>
	    {{ $url := urlFor "Lang.Index" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{ T "lang" 2}}</a></li>
        </ul>
        <form class="navbar-form navbar-left">
            <input id="search" type="text" class="form-control col-lg-8" placeholder="{{T "search"}}">
        </form>
        <ul class="nav navbar-nav navbar-right">
	    {{if .UserIsAdmin}}

	    {{ $url := urlFor "Trash.Get" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{ T "trash"}}</a></li>

	    {{ $url := urlFor "Debug.LastLog" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{ T "log" 2}}</a></li>
            {{end}}

	    {{ $url := urlFor "Prof.Get" }}
	    <li {{ if urlContain .rUrl $url }} class="active"{{end}}><a href={{ $url }}>{{.UserName}}</a></li>

{{range $k,$v := .Langs}}
<li{{if eq $k $.CurrentLang}} class="active"{{end}}><a href="/?lang={{$k}}">{{$v}}</a>
{{end}}
	    {{ $url := urlFor "Main.Logout" }}
            <li><a href="{{ $url }}">{{T "logout"}}</a></li>
        </ul>
    </div>
</div>
{{ else }}
	<center><h3>{{T "welcome"}}</h3></center>
{{ end }}
<!-- menu -->

  <div class="container page">
     {{.LayoutContent}}
     {{ if .gitlabToken }} 
     <button id="create_issue_btn" data-toggle="modal" data-target="#create_issue_modal" class="btn btn-fab btn-raised btn-primary">    
        <i class="mdi-action-bug-report"></i>
     </button>
     <div class="modal" id="create_issue_modal">
         <div class="modal-dialog">
             <div class="modal-content">
                 <div class="modal-header">
                     <button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button>
                     <h4 class="modal-title">{{T "issue_header"}}</h4> </div> <div class="modal-body">
                     <form method="post" id="issue_form">
                         <div class="form-group">
                             <input type="text" name="title" id="title" class="form-control" placeholder="{{T "issue_title"}}" required="true">
                         </div>
                         <div class="form-group">
                             <textarea class="form-control" placeholder="{{T "issue_description"}}" name="description" required="true"></textarea>
                         </div>
                         <div class="form-group">
                             <h4>Labels</h4>
                             <div id="issue_labels"></div>
                         </div> 
                         <input type="hidden" name="_xsrf"> 
                         <button type="submit" class="btn btn-primary">{{T "issue_submit"}}</button>
                     </form>   
                 </div>
             </div>
         </div>
     </div>
     {{end}}
  </div> <!-- /container -->

{{if .notice}}
<script>
$(document).ready(function(){
  show_notice("{{.notice}}")
})
</script>
{{end}}
</body>
</html>
