{{if gt .paginator.PageNums 1}}
<div style="position:relative;float:left;left:50%">
<div style="position:relative;float:right;right:50%">

<ul class="pagination pagination-sm">
    {{if .paginator.HasPrev}}
        <li><a href="{{.paginator.PageLinkFirst}}">{{T "page_first"}}</a></li>
        <li><a href="{{.paginator.PageLinkPrev}}">&lt;</a></li>
    {{else}}
        <li class="disabled"><a>{{T "page_first"}}</a></li>
        <li class="disabled"><a>&lt;</a></li>
    {{end}}
    {{range $index, $page := .paginator.Pages}}
        <li{{if $.paginator.IsActive .}} class="active"{{end}}>
            <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
        </li>
    {{end}}
    {{if .paginator.HasNext}}
        <li><a href="{{.paginator.PageLinkNext}}">&gt;</a></li>
        <li><a href="{{.paginator.PageLinkLast}}">{{T "page_last"}}</a></li>
    {{else}}
        <li class="disabled"><a>&gt;</a></li>
        <li class="disabled"><a>{{T "page_last"}}</a></li>
    {{end}}
</ul>
</div>
</div>
<div style="clear:both;"></div>
{{end}}
