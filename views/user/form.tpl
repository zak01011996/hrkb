    <h2>{{if .isEdit}}{{T "user_edit"}}{{else}}{{T "new_user"}}{{end}}</h2>
      {{if .isEdit}}
      {{ $url := urlFor "User.Get" ":id" .user.Id }}
      {{else}}
      {{ $url := urlFor "User.Add" }}
      {{end}}
	<form method="post" action="{{$.url}}">
	  {{.xsrfdata}}
	  <input type="hidden" name="id" value="{{if .isEdit}}{{.user.Id}}{{else}}0{{end}}">
	  
	  <div class="form-group{{if .errLogin}} has-error{{end}}">
	  	<label for="login" class="control-label">{{.errLogin}}</label>
	  	<input type="text" name="login" id="login" class="form-control" placeholder="{{T "login"}}" value="{{.user.Login}}"> 
	  </div>
          
          <div class="form-group{{if .errName}} has-error{{end}}">
	  	<label for="name" class="control-label">{{.errName}}</label>
	  	<input type="text" name="name" id="name" class="form-control" placeholder="{{T "name"}}" value="{{.user.Name}}"> 
	  </div>

	  <div class="form-group{{if .errEmail}} has-error{{end}}">
	  	<label for="mail" class="control-label">{{.errEmail}}</label>
	  	<input type="email" name="mail" id="mail" class="form-control" placeholder="Email" value="{{.user.Email}}"> 
	  </div>
            
          <div class="form-group{{if .errGToken}} has-error{{end}}">
	  	<label for="gtoken" class="control-label">{{.errGToken}}</label>
	  	<input type="text" name="gtoken" id="gtoken" class="form-control" placeholder="Auth token" value="{{.user.GToken.String}}"> 
	  </div>

	  <div class="form-group{{if .errPassword}} has-error{{end}}">
	  	<label for="pass" class="control-label">{{.errPassword}}</label>
	  	<input type="password" name="pass" id="pass" class="form-control" placeholder="{{T "password"}}"> 
	  </div>

	  <div class="form-group{{if .errPassConf}} has-error{{end}}">
	  	<label for="passc" class="control-label">{{.errPassConf}}</label>
	  	<input type="password" name="passc" id="passc" class="form-control" placeholder="{{T "passc"}}"> 
	  </div>

	  <div class="form-group{{if .errRole}} has-error{{end}}">
	  	<label for="role" class="control-label">{{.errRole}}</label>
	  	<select name="role" class="form-control" id="role">
			{{range $r := .roles}}
			{{if $.user}}
			<option value="{{$r.Id}}"{{if eq $r.Id $.user.Role}} selected{{end}}>{{$r.Name}}</option>
			{{else}}
			<option value="{{$r.Id}}">{{$r.Name}}</option>
			{{end}}
			{{end}}
	    </select>
	  </div>	

	  <button type="submit" class="btn btn-primary">{{if .isEdit}}{{T "update"}}{{else}}{{T "new"}}{{end}}</button>
	</form>
