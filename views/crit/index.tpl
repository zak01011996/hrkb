<h1>{{T "crit" 2}}</h1>

<div class="jumbotron">
	<a href="{{urlFor "Crit.Add"}}" class="btn btn-success pull-right">{{T "add_crit"}}</a>
	<div style="clear: both;"></div>
	<div class="panel-group" id="groups" role="tablist" aria-multiselectable="true">
		{{range $group := .groups}}
			<div class="panel panel-info">
    				<div class="panel-heading" role="tab" id="headingOne"
                                     data-toggle="collapse"
				     data-parent="#groups"
                                     style="cursor: pointer;" 
				     href="#collapse{{ $group.DepId }}">
     		 			<h4 class="panel-title">
        				<span><b>{{ $group.Department }}</b></span>
      					</h4>
    				</div>
    				<div id="collapse{{ $group.DepId }}"
				     class="panel-collapse collapse" 
				     role="tabpanel" aria-labelledby="headingOne">
      					<div class="panel-body">
						{{ if $group.Criterias }}
						<table class="table table-striped">
							{{range $criteria := $group.Criterias}}
							<tr>
								<td>{{ $criteria.Title }}</td>
								<td class="text-right">
									<a href="{{urlFor "Crit.Get" ":id" $criteria.Id}}" class="btn-sm btn-success" style="margin-right: 10px;">
										<i class="mdi-content-create"></i>
										<div class="ripple-wrapper"></div>
									</a>
									<a href="{{urlFor "Crit.Remove" ":id" $criteria.Id}}" class="btn-sm btn-danger b__remove">
										<i class="mdi-action-delete"></i>
										<div class="ripple-wrapper"></div>
									</a>	
								</td>
							</tr>
							{{end}}	
						</table>
						{{ else }}
							<h4 style="text-align: center;">{{T "no_crit_in_dep"}}</h4>
						{{end}}						
					</div>
				</div>
			</div>
		{{end}}
	</div>
</div>
