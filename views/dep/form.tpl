        <h2>{{if .isEdit}}{{T "edit_dep"}}{{else}}{{T "add_dep"}}{{end}}</h2>
      {{if .isEdit}}
      {{ $url := urlFor "Dep.Get" ":id" .dep.Id }}
      {{else}}
      {{ $url := urlFor "Dep.Add" }}
      {{end}}
	<form method="post" action="{{$.url}}">
	  {{.xsrfdata}}
	  <input type="hidden" name="id" value="{{if .isEdit}}{{.dep.Id}}{{else}}0{{end}}">
	  <div class="form-group{{if .errTitle}} has-error{{end}}">
	  	<label for="title" class="control-label">{{.errTitle}}</label>
	  	<input type="text" name="title" id="title" class="form-control" placeholder="{{T "title"}}" value="{{.dep.Title}}"> 
	  </div>	
	  <button type="submit" class="btn btn-primary">{{if .isEdit}}{{T "update"}}{{else}}{{T "new"}}{{end}}</button>
	</form>
