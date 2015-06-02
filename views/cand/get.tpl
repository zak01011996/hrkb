<h2 class="page_title">{{T "cand_card"}}</h2>
<div class="row">
    <div class="col-xs-12 col-md-9 center-block">
        <div class="panel panel-default cand_panel">
            <div class="panel-body">
                <div class="row">
                    <div class="col-md-3 col-sm-12">
                        <img src="/{{if .cand.Img.Valid}}{{.cand.Img.String}}{{else}}static/img/noavatar.png{{end}}" alt="Img" class="img-responsive">    
                    </div>

                    <div class="col-md-9 col-sm-12">
                            <table class="table table-striped">
                                    <tr class="align-top">
                                            <th>{{T "name"}}</th><td>{{.cand.Name}}</td>
                                    </tr>
                                    <tr>
                                            <th>{{T "last_name"}}</th><td>{{.cand.LName}}</td>
                                    </tr>
                                    <tr>
                                            <th>{{T "fath_name"}}</th><td>{{.cand.FName}}</td>
                                    </tr>
                                    <tr>
                                            <th>{{T "phone"}}</th><td>{{.cand.Phone}}</td>
                                    </tr>
                                    <tr>
                                            <th>{{T "email"}}</th><td>{{.cand.Email}}</td>
                                    </tr>

                                    <tr>
                                            <th>{{T "adr"}}</th><td>{{.cand.Address}}</td>
                                    </tr>
                                    <tr>
                                            <th>{{T "mar"}}</th><td>{{if .cand.Married}}{{T "married"}}{{else}}{{T "not_married"}}{{end}}</td>
                                    </tr>
                                    <tr>
                                            <th>{{T "dep" 1}}</th><td>{{index .deps .cand.Dep}}</td>
                                    </tr>
                                    <tr>
                                            <th>{{T "cost"}}</th><td>{{floatStr .cand.Salary}} {{.cand.Currency}}</td>
                                    </tr>

                            </table>
                    </div>
                </div>
            </div>
        </div>
        {{ if .cand.Note }}
        <div class="panel panel-default cand_panel">
            <div class="panel-body">
                <p style="font-weight: bold;">{{T "desc"}}</p>
                {{.cand.Note.String}}
            </div>
        </div>
        {{ end }}
        
        <div class="panel panel-default cand_panel">
            <div class="panel-body">
                <p style="font-weight: bold;">{{T "contact" 2}}</p>
                <table class="table table-striped" id="contacts">
        		<tbody>
        		{{range $c := .contacts}}
        			<tr>
        				<td>{{$c.Name}}</td>
        				<td>{{$c.Value}}</td>
        				<td class="text-right">
        					<a data-href="{{ urlFor "Contact.Remove" ":id" $.cand.Id ":cid" $c.Id }}" class="btn btn-danger delcont" title="{{T "del"}}">
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
                             <div class="form-group col-md-5">
                                 <label for="criteria" class="control-label">{{T "contact" 1}}</label>
                                 <div class="form-control-wrapper">
        			     <input type="text" id="contact_email" name="name" style="width:100%;" class="form-control" placeholder="Ex: mail">
                                     <span class="material-input"></span>
                                 </div>
                             </div>
                             <div class="form-group col-md-5">
                                 <label for="contact_val" class="control-label">{{T "val"}}</label>
                                 <div class="form-control-wrapper">
        	                     <input type="text" name="value" id="contact_val" class="form-control" style="width: 100%;" placeholder="Ex: example@mail.com">
                                     <span class="material-input"></span>
                                 </div>
                             </div>
               		     <div class="form-group col-md-2">
                                <button type="button" class="btn btn-primary" id="add_cont">{{T "add"}}</button>
                             </div> 
        		</div>
        	</form>
 
            </div>
        </div>
       
        <div class="panel panel-default cand_panel">
            <div class="panel-body">
                <p style="font-weight: bold;">
                    <span class="text-left">{{T "files"}}</span>
                    <a id="show_dropzone" style="float: right; cursor: pointer;">{{T "upload"}}</a>
                </p>
                   <div class="row" id="download">	
	            	{{range $f := .files}}
	            		<div class="file">
	            			<a href="{{ urlFor "Download.Get" ":id" $f.Id }}" target="_blank" title="{{$f.Name}}">
	            				<span class="mdi-editor-insert-drive-file"></span>
	            				{{$f.Name}} ({{kb $f.Size}}kb)
	            			</a>
	            		</div>
	            	{{end}}
	            </div>
	            <div id="dropzone" class="row" style="display: none;">
	            	<form action="{{ urlFor "Upload.Index" ":id" .cand.Id }}"
                                 class="dropzone"
                                 id="my-awesome-dropzone"
                                 enctype="multipart/form-data">
                             	{{.xsrfdata}}
                           </form>
	            </div>


            </div>
        </div>
        
        <div class="panel panel-default cand_panel">
                <div class="panel-body">
                    <p style="font-weight: bold;">
                        <span class="text-left">{{T "rat" 2}}</span>
                        <a id="show_less" style="float: right; cursor: pointer;">{{T "rat_g"}}</a>
                        <a id="show_detailed" style="float: right; cursor: pointer; padding-right: 15px;">{{T "rat_d"}}</a>
                        <a id="show_my_ratings" style="float: right; cursor: pointer; padding-right: 15px;">{{T "rat_my"}}</a>
                    </p>
                   
	            <div class="panel-group" id="ratings_table" role="tablist" aria-multiselectable="true"></div> 
	            <div class="panel-group" id="ratings_detailed_table" style="display: none;" role="tablist" aria-multiselectable="true"></div> 
	            <div class="panel-group" id="ratings_my_table" style="display: none;" role="tablist" aria-multiselectable="true"></div> 

                    <p style="font-weight: bold;">{{T "set_rat"}}</p>
                    <form id="rating_add_form" class="form-inline" style="width: 100%;" method="POST">
                        {{.xsrfdata}}
                        <div class="form-group col-md-8">
                            <label for="criteria" class="control-label">{{T "crit_choose"}}</label>
                            <div class="form-control-wrapper">
                                <select name="crit" class="form-control" style="width: 100%;" id="criteria">
				{{range $cg := .crits}}
				  <optgroup label="{{$cg.Department}}">
				{{range $c := $cg.Criterias}}
			            <option value="{{$c.Id}}">{{$c.Title}}</option>
				{{end}}
				  </optgroup>
				{{end}}
                                </select>
                                <span class="material-input"></span>
                            </div>
                        </div>
                        <div class="form-group col-md-2">
                            <label for="rating" class="control-label">{{T "mark"}}</label>
                            <div class="form-control-wrapper">
                                <select name="value" class="form-control" style="width: 100%;" id="rating">
                                    <option value="1">1</option>
                                    <option value="2">2</option>
                                    <option value="3">3</option>
                                    <option value="4">4</option>
                                    <option value="5" selected>5</option>
                                </select>
                                <span class="material-input"></span>
                            </div>
                        </div>
                        <input type="hidden" name="cand" value="{{.cand.Id}}">
                        <div class="form-group col-md-2">
                            <button type="submit" class="btn btn-primary">{{T "add"}}</button>
                        </div> 
                    </form> 
                </div>
        </div>
       
        <div class="panel panel-default cand_panel">
                <div class="panel-body">
                    <p style="font-weight: bold;">{{T "comments"}}</p>

                    {{ range $c := .Comments }}

<div class="comment">
   <span class="name">{{$c.Author}}</span>
   <span class="date">{{$c.Date}}</span>
   {{if eq $.UserId $c.User}}
   <span class="glyphicon glyphicon-pencil edit"></span>
   {{end}}
   {{if $.UserIsAdmin}}
   <a href="{{ urlFor "Comments.Remove" ":id" $c.Id }}" class="glyphicon glyphicon-remove remove"></a>
   {{end}}
   <pre class="text">{{$c.Comment | htmlunquote}}</pre>

   <form action="{{ urlFor "Comments.Edit" ":id" $c.Id }}" method="post">
     <div class="form-group has-error"><label class="control-label"></label></div>
     {{$.xsrfdata}}
     <div class="form-group">
      <textarea name="text" rows="4" class="form-control">{{$c.Comment | htmlunquote}}</textarea>
     </div>
     <div class="form-group">
       <button type="submit" class="btn btn-primary">{{T "save"}}</button> <button type="button" class="btn btn-primary">{{T "cancel"}}</button>
     </div>
   </form>
</div>

                    {{end}}

                    <form action="{{ urlFor "Comments.Get" ":id" .cand.Id }}" method="post" id="add-comment-form"{{if $.UserIsAdmin}} data-ad="1"{{end}}>
                     <div class="form-group has-error"><label class="control-label" for="text"></label></div>
                     {{.xsrfdata}}
                     <textarea name="text" id="text" class="form-control"></textarea>
                     <button type="submit" class="btn btn-primary">{{T "add_comment"}}</button>
                    </form>

                </div>
        </div>
    </div>
</div>
<script type="text/javascript">
    $(function() {
        Ratings.fillTables({{.cand.Id}});
    });
</script>
