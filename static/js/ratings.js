var Ratings = {}

$(function() {
    Ratings.mainURL = '/adm/ratings/';
    Ratings.detRatings = '/adm/ratings/detailed/';
    Ratings.myRatings = '/adm/ratings/my/';
    Ratings.ratingsForm = $('#rating_add_form');
    Ratings.grpBlock = $('#ratings_table');
    Ratings.detBlock = $('#ratings_detailed_table');
    Ratings.userBlock = $('#ratings_my_table');
    Ratings.btnDetailed = $('#show_detailed');
    Ratings.btnGrouped = $('#show_less');
    Ratings.btnMy = $('#show_my_ratings');

    Ratings.showNotice = function(msg) {
        $.snackbar({
            content: msg,
            style: "snackbar",
            timeout: 5000
        });
    };

    Ratings.getData = function(url, cb) {
        $.ajax({
            url: url,
            method: "GET",
            success: function(data) {
                if (data != null) {
                    cb(data.Ratings);
                }
            },
            error: function(r) {
                console.log(r);
            }
        });

    };

    Ratings.getPanelInfo = function(bodyId, title) {
        var panelInfo = $('<div/>', {
            'class': "panel panel-info"
        });

        var panelHeading = $('<div/>', {
            'class': "panel-heading",
            'role': "tab",
            'data-toggle': "collapse",
            'style': "cursor: pointer;",
            'href': "#" + bodyId
        });

        var titleHeader = $('<h4/>', {
            'class': "panel-title",
            'html': "<span><b>" + title + "</b></span>"
        });
        return panelInfo.append(panelHeading.append(titleHeader));
    };

    Ratings.getPanelBody = function(bodyId) {

        var panelBody = $('<div/>', {
            'id': bodyId,
            'class': "panel-collapse collapse",
            'role': "tabpanel",
        });

        return panelBody;
    };

    Ratings.getTableRow = function(criteria, rating, author, tds) {
        var row = $('<tr/>');
        var critTd = $('<td/>', {
            'text': criteria
        });
        row.append(critTd);

        if (author != "") {
            var authorTd = $('<td/>', {
                'text': author
            });
            row.append(authorTd);
        }

        var ratTd = $('<td/>', {
            'text': rating,
            'class': 'text-center text-rating'
        });
        row.append(ratTd);

        $.each(tds, function(i, td) {
            row.append(td);
        });

        return row;
    }

    Ratings.getDeps = function(data) {
        var currentDep = "";
        var deps = [];

        $.each(data, function(i, rating) {
            if (currentDep != rating.Dep) {
                deps.push(rating.Dep);
                currentDep = rating.Dep;
            }
        });

        return deps;
    };

    Ratings.fillGroupedRatings = function(data) {
        var self = this;
        self.grpBlock.empty();
        var deps = self.getDeps(data);

        $.each(deps, function(i, dep) {
            var ratings = _.where(data, {
                Dep: dep
            });

            var panelInfo = self.getPanelInfo("grouped_" + i, dep);
            var panelBody = self.getPanelBody("grouped_" + i);
            var grpTable = $('<table/>', {
                'class': 'table'
            });

            $.each(ratings, function(k, rating) {
                grpTable.append(self.getTableRow(rating.Criteria, rating.Rating, "", []));
            });

            self.grpBlock.append(panelInfo.append(panelBody.append(grpTable)));
        });
    };

    Ratings.fillDetailedRatings = function(data) {
        var self = this;
        self.detBlock.empty();
        var deps = self.getDeps(data);

        $.each(deps, function(i, dep) {
            var ratings = _.where(data, {
                Dep: dep
            });

            var panelInfo = self.getPanelInfo("detailed_" + i, dep);
            var panelBody = self.getPanelBody("detailed_" + i);
            var grpTable = $('<table/>', {
                'class': 'table'
            });

            $.each(ratings, function(k, rating) {
                grpTable.append(self.getTableRow(rating.Criteria, rating.Rating, rating.Author, []));
            });

            self.detBlock.append(panelInfo.append(panelBody.append(grpTable)));
        });
    };

    Ratings.fillMyRatings = function(data) {
        var self = this;
        self.userBlock.empty();
        var deps = self.getDeps(data);

        $.each(deps, function(i, dep) {
            var ratings = _.where(data, {
                Dep: dep
            });

            var panelInfo = self.getPanelInfo("my_" + i, dep);
            var panelBody = self.getPanelBody("my_" + i);
            var grpTable = $('<table/>', {
                'class': 'table',
                'click': function(e) {
                    var elem = $(e.target);
                    if (elem.hasClass('rating-del')) {
                        self.deleteRating(elem.data('id'));
                    }
                }
            });

            $.each(ratings, function(k, rating) {
                var btnTd = $('<td/>', {
                    'class': 'text-right'
                });

                var delBtn = $('<a/>', {
                    'class': 'btn btn-danger btn-sm rating-del',
                    'data-id': rating.Id,
                    'text': 'DELETE'
                }); 
                
                btnTd.append(delBtn);
                grpTable.append(self.getTableRow(rating.Criteria, rating.Rating, "", [btnTd]));
            });

            self.userBlock.append(panelInfo.append(panelBody.append(grpTable)));
        });
    };

    Ratings.fillTables = function(candId) {
        var self = this;
        self.getData(this.mainURL + candId, function(data) {
            self.fillGroupedRatings(data);
        });

        self.getData(this.detRatings + candId, function(data) {
            self.fillDetailedRatings(data);
        });

        self.getData(this.myRatings + candId, function(data) {
            self.fillMyRatings(data);
        });
    };

    Ratings.deleteRating = function(ratingId) {
        var self = this;
        var candId = self.ratingsForm.find("[name=cand]").val();
        $.ajax({
            url: self.mainURL + ratingId + '/remove',
            method: "GET",
            success: function(data) {
                self.fillTables(candId);
                self.showNotice(data.message);
            },
            error: function(r) {
                console.log(r);
            }
        });
    };

    Ratings.ratingsForm.submit(function() {
        var candId = Ratings.ratingsForm.find("[name=cand]").val();
        $.ajax({
            url: Ratings.mainURL + "0",
            method: "POST",
            data: Ratings.ratingsForm.serialize(),
            success: function(r) {
                Ratings.showNotice(r.message);
                Ratings.fillTables(candId);
            },
            error: function(r) {
                console.log(r);
            }
        });

        return false;
    });

    Ratings.btnDetailed.click(function() {
        Ratings.grpBlock.hide();
        Ratings.userBlock.hide();
        Ratings.detBlock.show();
    });

    Ratings.btnGrouped.click(function() {
        Ratings.userBlock.hide();
        Ratings.detBlock.hide();
        Ratings.grpBlock.show();
    });

    Ratings.btnMy.click(function() {
        Ratings.detBlock.hide();
        Ratings.grpBlock.hide();
        Ratings.userBlock.show();
    });
});
