<h2>{{T "prof"}}</h2>
<form method="post" action="{{urlFor "Prof.Get"}}">
	{{.xsrfdata}}
	<div class="form-group{{if .errName}} has-error{{end}}">
	    <label for="name" class="control-label">{{.errName}}</label>
	    <input type="text" id="name" name="name" class="form-control" placeholder="{{T "name"}}" value="{{.Name | html}}">
	</div>
	<div class="form-group{{if .errEmail}} has-error{{end}}">
	    <label for="email" class="control-label">{{.errEmail}}</label>
	    <input type="text" id="email" name="email" class="form-control" placeholder="{{T "email"}}" value="{{.user.Email | html}}">
	</div>
       	<h4>{{T "changepass"}}</h4>

	<div class="form-group{{if .errPassword}} has-error{{end}}">
	  	<label for="pass" class="control-label">{{.errPassword}}</label>
	  	<input type="password" name="pass" id="pass" class="form-control" placeholder="{{T "password"}}"> 
	</div>

	<div class="form-group{{if .errPassConf}} has-error{{end}}">
	  	<label for="passc" class="control-label">{{.errPassConf}}</label>
	  	<input type="password" name="passc" id="passc" class="form-control" placeholder="{{T "passc"}}"> 
	</div>
	<div class="form-group">
		<label class="control-label">Notify</label>  <br>
		<div class="togglebutton">
                	<label>
				<input type="checkbox" id="notifyMail" name="notifyMail" value="true" {{if .user.NotifyByMail}} checked="true" {{end}}> 
				By Mail
               		</label>
              	</div>
		<div class="togglebutton">
                	<label>
				<input type="checkbox" id="notifyMail" name="notifyTelegram" value="true" {{if .user.NotifyByTelegram}} checked="true" {{end}}>
				By Telegram 
               		</label>
              	</div>
	</div>

	<button type="submit" class="btn btn-primary">{{T "save"}}</button>
</form>
