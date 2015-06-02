<div class="jumbotron">

	<h3>{{T "contact" 2}}</h3>
	<table class="table table-striped" id="contacts">
		<tbody>
		{{range $c := .contacts}}
			<tr>
				<td>{{$c.Name}}</td>
				<td>{{$c.Value}}</td>
				<td class="text-right">
					<a data-href="{{ urlFor "Contact.Remove" ":id" $.cand.Id ":cid" $c.Id }}" class="btn btn-danger delcont" title="Delete">
						<span class="glyphicon glyphicon-trash"></span>
					</a>
				</td>
			</tr>
		{{end}}
		</tbody>
	</table>

	<form action="{{ urlFor "Contact.Add" ":id" .cand.Id }}" class="form-inline" method="post">
		{{.xsrfdata}}
		<div class="row">
			<div class="col-xs-4">
				<input type="text" name="name" class="form-control" placeholder="mail">
			</div>
			<div class="col-xs-4">
				<input type="text" name="value" class="form-control" placeholder="example@mail.com">
			</div>
			<div class="col-xs-4">
				<input type="button" class="btn btn-primary" value="{{T "add"}}" id="add_cont">
			</div>
		</div>
	</form>
</div>
