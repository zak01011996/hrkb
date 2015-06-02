<div class="jumbotron">
    {{if .authErr }}
        <h4 style="text-align: center;">{{.authErr}}</h4>
    {{else}}
    <h2>{{T "issue_header"}}</h2>
        <form method="post" action="/adm/issue/create" id="issue_form">
        {{.xsrfdata}}
          
        <div class="form-group">
            <input type="text" name="title" id="title" class="form-control" placeholder="{{T "issue_title"}}">
        </div>
        <div class="form-group">
            <textarea class="form-control" placeholder="{{T "issue_description"}}" name="description" required="true"></textarea>
        </div>
        <div class="form-group">
            <h4>Labels</h4>
            {{ range $label := .labels }}
                <label><input type="checkbox" name="labels[]" value="{{$label.Name}}">{{$label.Name}}</label><br/>
            {{end}}
        </div> 
        
        <button type="submit" class="btn btn-primary">{{T "issue_submit"}}</button>
    </form>
    {{end}}
</div>
