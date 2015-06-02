<h2>{{if .crit}}{{T "edit_crit"}}{{else}}{{T "new_crit"}}{{end}}</h2>
<div class="jumbotron">
      <h3>{{T "fill"}}</h3>
      {{if .crit}}
      {{ $url := urlFor "Crit.Get" ":id" .crit.Id }}
      {{else}}
      {{ $url := urlFor "Crit.Add" }}
      {{end}}
      <form method="post" action="{{ $.url }}">
	  {{.xsrfdata}}
         <div class="form-group{{if .errTitle}} has-error{{end}}">
               <label for="title" class="control-label">{{.errTitle}}</label>
	       <input type="text" id="title" name="title" class="form-control" placeholder="{{T "title"}}" value="{{if .crit}}{{.crit.Title}}{{end}}">
         </div>

         <div class="form-group">
            <label for="dep" class="control-label">
                  {{T "dep" 1}}
            </label>
	    <select name="dep" id="dep" class="form-control">
               {{ range $dep := .deps }}
                  <option value="{{ $dep.Id }}" 
                     {{if $.crit}}
                        {{if eq $dep.Id $.crit.Dep}}
                           selected
                        {{end}}
                     {{end}}>
                    {{ $dep.Title }} 
                  </option>
               {{ end }}
	    </select>
         </div>

         <input type="hidden" name="id" value="{{if .crit}}{{.crit.Id}}{{end}}">

	 <a href="{{urlFor "Crit.Index"}}" type="submit" class="btn btn-warning">{{T "back"}}</a>
	 <button type="submit" class="btn btn-primary">
                  {{if .crit}}{{T "edit"}}{{else}}{{T "add"}}{{end}}
         </button>
      </form>
</div>
