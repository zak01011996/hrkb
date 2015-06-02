<div class="jumbotron">

	<h3>{{T "file" 2}}</h3>
	<div class="list-group">
		{{range $f := .files}}
    			<div class="list-group-item">
    			    <div class="row-action-primary">
				<a href="/{{$f.Url}}" target="blank">
			    	{{if isImg $f.Mime}}
					<img src="/{{$f.Url}}" class="circle" alt="icon">
				{{else}}
    			        	<i class="mdi-editor-insert-drive-file"></i>
				{{end}}
				</a>
    			    </div>
    			    <div class="row-content">
    			        <div class="action-secondary"><i class="mdi-material-info"></i></div>
    			        <h4 class="list-group-item-heading">{{$f.Name}}</h4>
    			        <p class="list-group-item-text">
					{{T "type"}}: {{$f.Mime}}, {{kb $f.Size}} KB 
					<span class="pull-right"><button class="btn btn-danger delfile" data-href="{{urlFor "Download.Remove" ":id" $f.Id }}">{{T "del"}}</button></span>
				</p>
    			    </div>
    			</div>
			<div class="list-group-separator"></div>
    		{{end}}
	</div>
	<div class="row">
		<form action="{{urlFor "Upload.Img" ":id" 1}}"
                      class="dropzone"
                      id="my-awesome-dropzone"
                      enctype="multipart/form-data">
                  	{{.xsrfdata}}
                </form>
	</div>

</div>
