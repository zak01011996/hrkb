<h1>{{T "cand" 2}}</h1>

<div class="jumbotron">

<form class="form-inline" method="get" action="{{ urlFor "Cand.Index" }}" id="cand_filter" style="margin-bottom:20px;">
  <div class="form-group">
    <label>{{T "dep_filter"}}</label>
    <select name="dep" class="form-control">
    <option value="0">{{T "dep_filter_all"}}</option>
    {{range $k,$v := .deps}}
    <option value="{{$k}}"{{if eq $k $.depFilter}} selected{{end}}>{{$v}}</option>
    {{end}}
    </select>
  </div>
</form>

	<a href="{{urlfor "Cand.Add"}}" class="btn btn-primary"><span class="glyphicon glyphicon-plus"></span>{{T "add_cand"}}</a>
	<div class="table-responsive">
	<table class="table table-striped table-candidates">
		<thead>
			<tr>
				<th></th>
				<th colspan="2">{{T "info"}}</th>
				<th>{{T "desc"}}</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
		{{range $c := .cands}}
			<tr>
				<td>
                                    <a href="{{ urlFor "Cand.Get" ":id" $c.Id }}">
                                    <img src="/{{if $c.Img.Valid}}{{$c.Img.String}}{{else}}static/img/noavatar.png{{end}}" 
                                         alt="Img"
                                         width=100>
                                    </a>
                                </td>
				<td colspan="2">
                                    <span class="candidate-info">
                                        <p>
                                            <span style="font-weight: bold;">{{T "name"}}: </span>{{$c.Name}}
                                        </p>
                                        <p>
                                            <span style="font-weight: bold;">{{T "last_name"}}: </span>{{$c.LName}}
                                        <p>
                                            <span style="font-weight: bold;">{{T "dep" 1}}: </span>{{index $.deps $c.Dep}}
                                        </p>
                                        <p>
                                            <span style="font-weight: bold;">{{T "phone"}}: </span>{{$c.Phone}}
                                        </p>
                                        <p>
                                            <span style="font-weight: bold;">{{T "email"}}: </span>{{$c.Email}}
                                        </p>

                                    </span>
                                </td>
                                <td style="max-width: 300px;">
                                    <span class="candidate-info">
                                        {{cutStr $c.Note.String 150}}
                                    </span>
                                </td>
				<td class="text-right">
					<a href="{{ urlFor "Cand.Get" ":id" $c.Id }}" class="btn btn-primary" title="{{T "view"}}">
						<span class="glyphicon glyphicon-eye-open"></span>
					</a>
					<a href="{{ urlFor "Cand.Edit" ":id" $c.Id }}" class="btn btn-warning" title="{{T "edit"}}">
						<span class="glyphicon glyphicon-pencil"></span>
					</a>
					<a href="{{ urlFor "Cand.Remove" ":id" $c.Id }}" class="btn btn-danger b__remove" title="{{T "del"}}">
						<span class="glyphicon glyphicon-trash"></span>
					</a>
				</td>
			</tr>
		{{end}}
		</tbody>
	</table>
	</div>
	{{template "pages.tpl" .}}
</div>
