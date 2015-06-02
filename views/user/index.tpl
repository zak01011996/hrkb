<h1>{{T "user" 2}}</h1>

<div class="jumbotron">
	<a href="{{urlFor "User.Add"}}" class="btn btn-primary"><span class="glyphicon glyphicon-plus"></span>{{T "add_user"}}</a>
	<div class="table-responsive">
		<table class="table table-striped">
			<thead>
				<tr>
					<th>ID</th>
					<th>Login</th>
					<th>Name</th>
					<th>Email</th>
					<th>Role</th>
					<th>Notify Mail</th>
					<th>Notify Telegram</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
			{{range $user := .users}}
				<tr>
					<td>{{$user.Id}}</td>
					<td><a href="{{urlFor "User.Get" ":id" $user.Id}}">{{$user.Login}}</a></td>
					<td><a href="{{urlFor "User.Get" ":id" $user.Id}}">{{$user.Name}}</a></td>
					<td>{{$user.Email}}</td>
					<td>{{index $.roles $user.Role}}</td>
					<td class="text-center">
						{{if $user.NotifyByMail}}
							<div class="icon-preview">
								<i class="mdi-toggle-check-box"></i>
							</div>
						{{else}}
							<div class="icon-preview">
								<i class="mdi-toggle-check-box-outline-blank"></i>
							</div>
						{{end}}
					</td>
					<td class="text-center">
						{{if $user.NotifyByTelegram}}
							<div class="icon-preview">
								<i class="mdi-toggle-check-box"></i>
							</div>
						{{else}}
							<div class="icon-preview">
								<i class="mdi-toggle-check-box-outline-blank"></i>
							</div>
						{{end}}
					</td>
					<td><a href="{{urlFor "User.Remove" ":id" $user.Id}}" class="btn btn-danger b__remove">{{T "delete"}}</a></td>                
				</tr>
			{{end}}
			</tbody>
		</table>
	</div>
</div>
