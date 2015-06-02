<div class="jumbotron">
	<div class="table-responsive">
	<table class="table table-striped">
		<thead>
			<tr>
				<th width="15%">{{T "date"}}</th>
				<th width="15%">{{T "user" 1}}</th>
				<th width="15%">{{T "cand" 1}}</th>
				<th>{{T "text"}}</th>
                                <th width="20%"></th>
			</tr>
		</thead>
		<tbody>
                {{range $c := .comments}}
                <tr>
                  <td>{{$c.Dt}}</td>
                  <td>{{$c.UserName}}</td>
                  <td>{{$c.CandName}}</td>
                  <td>{{$c.Text | htmlunquote}}</td>
		  <td><a href="{{urlFor "Trash.Restore" ":type" "comments" ":id" $c.Id}}" class="btn btn-success b__restore">{{T "restore"}}</a></td> 
                </tr>
                {{end}}
		</tbody>
	</table>
</div>
</div>
