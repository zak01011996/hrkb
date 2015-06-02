<h1>{{T "lang" 2}}</h1>

<div class="jumbotron">
	<table class="table table-striped">
          <thead>
            <tr>
                <th>{{T "title"}}</th>
                <th></th>
            </tr>
          </thead>
          <tbody>
          {{range $l := .langs}}
          <tr{{if $l.IsDefault}} class="danger"{{end}}>
            <td>{{index $.langsSup $l.Code}}</td>
            <td>
                {{if not $l.IsDefault}}
		<a href="{{urlFor "Lang.Default" ":id" $l.Id}}" class="btn btn-primary">{{T "lang_setdefault"}}</a>
                {{end}}
		<div class="btn btn-primary lang-upload" data-id="{{$l.Id}}"><span class="glyphicon glyphicon-upload"></span> {{T "lang_upload"}}</div>
		<a href="{{urlFor "Lang.Download" ":id" $l.Id}}" class="btn btn-primary"><span class="glyphicon glyphicon-download"></span> {{T "lang_download"}}</a>
		<a href="{{urlFor "Lang.Remove" ":id" $l.Id}}" class="btn btn-danger b__remove" title="{{T "del"}}"><span class="glyphicon glyphicon-trash"></span></a>
            </td>
          </tr>
          {{end}}
          </tbody>
        </table>
        <div class="dropzone-previews" id="langs-previews"></div>

        <h3>{{T "lang_add"}}</h3>

        <form action="{{urlFor "Lang.Add"}}" method="post" enctype="multipart/form-data" id="lang-add-form">
        {{.xsrfdata}}

	  <div class="form-group{{if .errCode}} has-error{{end}}">
	    <label for="code" class="control-label">{{.errCode}}</label>
	    <select name="code" class="form-control" id="code">
              {{range $k,$v := .langsSup}}
              <option value="{{$k}}"{{if eq $.code $k}} selected{{end}}>{{$v}}</option>
              {{end}}
	    </select>
	  </div>

          <div class="form-group{{if .errFile}} has-error{{end}}">

	    <label for="tfile" class="control-label">{{.errFile}}</label>

            <label>
	    <a class="btn btn-primary">{{T "lang_chosefile"}}</a>
	    <input type="file" name="file" id="tfile" class="form-control">
            </label><span class="l__filename"></span>

          </div>

	  <button type="submit" class="btn btn-primary">{{T "add"}}</button>
        </form>
</div>
