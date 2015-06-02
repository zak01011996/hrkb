<h2>{{T "add_cand"}}</h2>
<div class="col-md-6">

{{if .flash.error}}
<div class="alert alert-danger" role="alert">
	{{.flash.error}}
</div>
{{end}}

<p>Fields with <span class="required">* </span> are required</p>

{{ $url := urlFor "Cand.Add" }}
<form method="post" action="{{ $url }}" enctype="multipart/form-data">
  {{.xsrfdata}}

  <div class="form-group{{if .errName}} has-error{{end}}">
	<label for="name" class="control-label">
	{{if .errName}}
		{{.errName}}
	{{else}}
	<span class="required">* </span>{{T "name"}}
	{{end}}
	</label>
	<input type="text" name="name" id="name" class="form-control" placeholder="{{T "name"}}" value="{{.cand.Name}}">
  </div>

 <div class="form-group{{if .ErrImg}} has-error{{end}}">
	<label for="img" class="control-label">
	{{if .ErrImg}}
		{{.ErrImg}}
	{{else}}
		{{T "img"}} <br>
		<a class="btn btn-primary">{{T "choose_img"}}</a>
	{{end}}
	</label>
	<input type="file" name="img" id="img" class="form-control">
  </div>

  <div class="form-group{{if .errLName}} has-error{{end}}">
	<label for="lname" class="control-label">
	{{if .errLName}}
		{{.errLName}}
	{{else}}
		<span class="required">* </span>{{T "last_name"}}
	{{end}}
	</label>
	<input type="text" name="lname" id="lname" class="form-control" placeholder="{{T "last_name"}}" value="{{.cand.LName}}">
  </div>	

  <div class="form-group{{if .errFName}} has-error{{end}}">
	<label for="fname" class="control-label">
	{{if .errFName}}
		{{.errFName}}
	{{else}}
		<span class="required">* </span>{{T "fath_name"}}
	{{end}}
	</label>
	<input type="text" name="fname" id="fname" class="form-control" placeholder="{{T "fath_name"}}" value="{{.cand.FName}}">
  </div>	
  <div class="form-group{{if .errPhone}} has-error{{end}}">
	<label for="phone" class="control-label">
	{{if .errPhone}}
		{{.errPhone}}
	{{else}}
	    {{T "phone"}}	
	{{end}}
	</label>
	<input type="number" name="phone" id="phone" class="form-control" placeholder="{{T "phone"}}" value="{{.cand.Phone}}">
  </div>	
  <div class="form-group{{if .errEmail}} has-error{{end}}">
	<label for="email" class="control-label">
	{{if .errEmail}}
		{{.errEmail}}
	{{else}}
	        <span class="required">* </span>{{T "email"}}	
	{{end}}
	</label>
	<input type="email" name="email" id="email" class="form-control" placeholder="{{T "email"}}" value="{{.cand.Email}}">
  </div>
  <div class="form-group">
	<label for="note" class="control-label">
		{{T "notes"}}	
	</label>
	<textarea id="note" name="note" class="form-control" placeholder="{{T "notes"}}">{{.cand.Note.String}}</textarea>
  </div>

  <div class="form-group{{if .errAddress}} has-error{{end}}">
	<label for="addr" class="control-label">
	{{if .errAddress}}
		{{.errAddress}}
	{{else}}
		<span class="required">* </span>{{T "adr"}}	
	{{end}}
	</label>
	<textarea id="addr" name="addr" class="form-control" placeholder="{{T "adr"}}">{{.cand.Address}}</textarea>
  </div>

  <div class="form-group{{if .errMarried}} has-error{{end}}">
	<label for="married" class="control-label">
	{{if .errMarried}}
		{{.errMarried}}
	{{else}}
		{{T "married"}}
	{{end}}
	</label> <br/>
	<input type="checkbox" name="married" id="married" value="true">
  </div>

   <div class="form-group{{if .errDep}} has-error{{end}}">
	<label for="depId" class="control-label">
	{{if .errDep}}
		{{.errDep}}
	{{else}}
		{{T "dep" 1}}
	{{end}}
	</label>
	<select id="depId" name="depId" class="form-control">
	{{range $dep := .deps}}
		<option value="{{$dep.Id}}">{{$dep.Title}}</option>
	{{end}}
	</select>
  </div>  

  <div class="form-group{{if .errSalary}} has-error{{end}}">
	<label for="salary" class="control-label">
	{{if .errSalary}}
		{{.errSalary}}
	{{else}}
		<span class="required">* </span>{{T "cost"}}	
	{{end}}
	</label>
	<input type="number" name="salary" id="salary" class="form-control" placeholder="{{T "cost"}}" value="{{floatStr .cand.Salary}}">
  </div>

  <div class="form-group{{if .errCurrency}} has-error{{end}}">
	<label for="currency" class="control-label">
	{{if .errCurrency}}
		{{.errCurrency}}
	{{else}}
		<span class="required">* </span>{{T "currency"}}	
	{{end}}
	</label>
	<input type="text" name="currency" id="currency" class="form-control" placeholder="{{T "currency"}}" value="{{.cand.Currency}}">
  </div>

  <button type="submit" class="btn btn-primary">{{T "new"}}</button>
</form>
</div>

