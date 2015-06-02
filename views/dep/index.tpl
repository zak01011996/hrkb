<h1>{{T "dep" 2}}</h1>

<div class="jumbotron">
	<a href="{{urlFor "Dep.Add"}}" class="btn btn-primary">{{T "add_dep"}}</a>
	<table class="table table-striped">
		<thead>
			<tr>
				<th>{{T "title"}}</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
		{{range $d := .deps}}
			<tr>
				<td><a href="{{urlFor "Dep.Get" ":id" $d.Id}}">{{$d.Title}}</a></td>
				<td><a href="{{urlFor "Dep.Remove" ":id" $d.Id}}" class="btn btn-danger b__remove">{{T "del"}}</a></td>
			</tr>
		{{end}}
		</tbody>
	</table>
</div>
