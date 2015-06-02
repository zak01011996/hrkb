<div role="tabpanel">

  <!-- Nav tabs -->
  <ul class="nav nav-tabs jstabs" role="tablist">
    <li role="presentation" class="active"><a href="#users" aria-controls="users" role="tab" data-toggle="tab">{{T "user" 2}}</a></li>
    <li role="presentation"><a href="#deps" aria-controls="deps" role="tab" data-toggle="tab">{{T "dep" 2}}</a></li>
    <li role="presentation"><a href="#crit" aria-controls="crit" role="tab" data-toggle="tab">{{T "crit" 2}}</a></li>
    <li role="presentation"><a href="#cand" aria-controls="cand" role="tab" data-toggle="tab">{{T "cand" 2}}</a></li>
    <li role="presentation"><a href="#comment" aria-controls="comment" role="tab" data-toggle="tab">{{T "comments" 2}}</a></li>
  </ul>

  <!-- Tab panes -->
<div class="tab-content">
  <div role="tabpanel" class="tab-pane fade in active" id="users">
   {{template "trash/users.tpl" .}}
  </div>
  <div role="tabpanel" class="tab-pane fade" id="deps">
   {{template "trash/deps.tpl" .}}
  </div>
  <div role="tabpanel" class="tab-pane fade" id="crit">
   {{template "trash/crits.tpl" .}}
  </div>
  <div role="tabpanel" class="tab-pane fade" id="cand">
   {{template "trash/cands.tpl" .}}
  </div>
  <div role="tabpanel" class="tab-pane fade" id="comment">
   {{template "trash/comments.tpl" .}}
  </div>
</div>

</div>
