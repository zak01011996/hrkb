<div class="jumbotron">
	<div class="table-responsive">
	<table class="table table-striped">
		<thead>
			<tr>
				<th width="40%">{{T "login"}}</th>
				<th width="40%">{{T "name"}}</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
		{{range $user := .users}}
			<tr>
                                <td>{{$user.Login}}</td>
				<td>{{$user.Name}}</td>
				<td><a href="{{urlFor "Trash.Restore" ":type" "users" ":id" $user.Id}}" class="btn btn-success b__restore">{{T "restore"}}</a></td> 
			</tr>
		{{end}}
		</tbody>
	</table>
</div>
</div>
