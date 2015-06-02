<div class="jumbotron">

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
                                    <img src="/{{if $c.Img.Valid}}{{$c.Img.String}}{{else}}static/img/noavatar.png{{end}}" 
                                         alt="Img"
                                         width=100>
                                </td>
				<td colspan="2">
                                    <span class="candidate-info">
                                        <p>
                                            <span style="font-weight: bold;">{{T "name"}}: </span>{{$c.Name}}
                                        </p>
                                        <p>
                                            <span style="font-weight: bold;">{{T "last_name"}}: </span>{{$c.LName}}
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
<a href="{{urlFor "Trash.Restore" ":type" "cands" ":id" $c.Id}}" class="btn btn-success b__restore">{{T "restore"}}</a>
				</td>
			</tr>
		{{end}}
		</tbody>
	</table>
	</div>
</div>
