<div class="jumbotron">
     	<div class="table-responsive">
	<table class="table table-striped">
		<thead>
			<tr>
				<th width="80%">{{T "title"}}</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
		{{range $d := .deps}}
			<tr>
                                <td>{{$d.Title}}</td>
				<td><a href="{{urlFor "Trash.Restore" ":type" "deps" ":id" $d.Id}}" class="btn btn-success b__restore">{{T "restore"}}</a></td>
			</tr>
		{{end}}
		</tbody>
	</table>
</div>
</div>
