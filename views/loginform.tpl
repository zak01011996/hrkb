<div id="auth" class="col-xs-12 col-md-4 center-block">
	<form method="post" action="{{urlFor "Main.Get"}}">
	  {{.xsrfdata}}
	  <div class="form-group{{if .errLogin}} has-error{{end}}">
	    <label for="login" class="control-label">{{.errLogin}}</label>
	    <input type="text" id="login" name="login" class="form-control" placeholder="{{T "login"}}" value="{{.login | html}}">
	  </div>
	  <div class="form-group{{if .errPassword}} has-error{{end}}">
	    <label for="password" class="control-label">{{.errPassword}}</label>
	    <input type="password" name="password" id="password" class="form-control" placeholder="{{T "password"}}">
	  </div>
	  <button type="submit" class="btn btn-primary">{{T "site_enter"}}</button>
	</form>
</div>
