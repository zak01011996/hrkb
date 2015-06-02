<ul class="list-group">
{{ range $log := .logs }}
   <li class="list-group-item"> {{ $log }} </li>
{{end}}
</ul>