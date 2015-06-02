// TODO: Refactor JS code with latest code standarts
xsrf = base64_decode($.cookie("_xsrf").split("|")[0]);
//xsrf field you can use it to send forms by ajax
$(document).ready(function() {

    $("#cand_filter select").change(function() {
        $(this).parents('form').submit();
    });

    $('.jstabs').click(function(e) {
        e.preventDefault()
        $(this).tab('show')
    });

    $(".delfile").click(function(e) {
        elem = $(this)
        $.ajax({
            method: "post",
            url: elem.data("href"),
            data: {
                _xsrf: xsrf
            },
            error: function(data) {
                show_notice(data.responseJSON.message)
            },
            success: function(data) {
                elem.parents(".list-group-item").remove()
                show_notice(data.message)
            },
        })
    });

    Dropzone.options.myAwesomeDropzone = {
        init: function() {
            this.on("success", function(file, data) {
                $("#download").append("<div class='file'><a href='/adm/download/" + data.Id + "' target='_blank'>" +
                    "<span class='mdi-editor-insert-drive-file'></span>" +
                    data.Name + " " + data.Size + " KB" +
                    "</a></div>")
            });
        }
    };

    $("#add_cont").click(function(e) {
        var btn = $(this);
        var data = $(this).parents("form").serialize();
        var id = $("[name=cand]").val();

        var name = btn.parents("form").find("[name=name]");
        var nameGroup = name.parents(".form-group")

        var value = btn.parents("form").find("[name=value]");
        var valueGroup = value.parents(".form-group")

        if (name.val().length < 3 || value.val().length < 3) {

            if (!nameGroup.hasClass("has-error")) {
                nameGroup.addClass("has-error")
            }

            if (!valueGroup.hasClass("has-error")) {
                valueGroup.addClass("has-error")
            }

            show_notice("Contact Type & Value must be at least 3 characters")
            return
        }

        $.ajax({
            method: "post",
            url: "/adm/candidates/" + id + "/contacts/add",
            data: btn.parents("form").serialize(),
            success: function(data) {

                nameGroup.removeClass("has-error")
                valueGroup.removeClass("has-error")

                $("#contacts>tbody").append("<tr><td>" +
                    name.val() + "</td><td>" +
                    value.val() + "</td>" +
                    "<td class='text-right'><a class='btn btn-danger delcont' data-href='/adm/candidates/" + id + "/contacts/" + data.id + "/remove'><span class='glyphicon glyphicon-trash'></span></a></td></tr>"
                );

                value.val("")
                name.val("")

            }
        })

    })

    $("#contacts").click(function(e) {
        var elem = $(e.target)

        if (elem.hasClass("delcont") || elem.parent().hasClass("delcont")) {
            $.ajax({
                url: elem.data("href"),
                success: function(data) {
                    elem.parents("tr").remove()
                    show_notice(data.message)
                }
            })
        }
    })

    $("td a.b__remove,td a.b__restore").click(function() {

        var t = $(this)
        if (!t.hasClass("active")) {
            t.addClass("active")
            $.ajax({
                url: t.attr("href"),
                success: function(r) {

                    if (r.success) {

                        t.parent().parent().remove()
                    }

                    show_notice(r.message)
                    t.removeClass("active")
                },
                error: function() {
                    t.removeClass("active")
                }
            })
        }

        return false
    });

    $('#add-comment-form').submit(function() {
        var t = $(this);

        if ($("button[type='submit']", t).hasClass("active")) {
            return false;
        }
        $("button[type='submit']", t).addClass("active");

        $.ajax({
            url: t.attr('action'),
            type: 'POST',
            data: t.serialize(),
            success: function(r) {
                $("button[type='submit']", t).removeClass("active");
                if (r.Msg == "Ok") {
                    $('#add-comment-form .has-error label').text("");

                    var s = '<div class="comment"><span class="name">' + r.Author + '</span> <span class="date">' + r.Date + '</span> <span class="glyphicon glyphicon-pencil edit"></span>';
                    if ($('#add-comment-form').attr('data-ad') == '1') {
                        s += ' <a href="/adm/comments/' + r.Id + '/remove" class="glyphicon glyphicon-remove remove"></a>';
                    }
                    s += '<pre class="text">' + r.Comment + '</pre>';
                    s += '<form action="/adm/comments/' + r.Id + '/edit" method="post">';
                    s += '<div class="form-group has-error"><label class="control-label"></label></div>';
                    s += '<input type="hidden" name="_xsrf" value="' + xsrf + '">';
                    s += '<div class="form-group"><textarea name="text" rows="4" class="form-control">' + r.Comment + '</textarea></div>';
                    s += '<div class="form-group">';
                    s += '<button type="submit" class="btn btn-primary">' + t_save + '</button> <button type="button" class="btn btn-primary">' + t_cancel + '</button>';
                    s += '</div></form></div>';

                    $(s).insertBefore('#add-comment-form').fadeIn('normal');
                    $('#add-comment-form textarea').val('');

                } else {
                    $('#add-comment-form .has-error label').text(r.Msg);
                }
            },
            error: function() {
                $("button[type='submit']", t).removeClass("active");
            }
        });

        return false;
    });

    $('#search').autocomplete({
        autoFocus: true,
        source: function(request, response) {
            $.ajax({
                url: '/adm/search',
                type: 'GET',
                data: {
                    q: request.term
                },
                success: function(data) {
                    response(data)
                }
            });
        },
        select: function(event, ui) {
            location.href = "/adm/candidates/" + ui.item.Id;
        },
        open: function(event, ui) {
            $('#ui-id-1.ui-autocomplete').css('margin-top', '-2px');
        },
        close: function(event, ui) {
            $('#ui-id-1.ui-autocomplete').css('margin-top', '0');
        }
    }).autocomplete("instance")._renderItem = function(ul, item) {
        return $('<li class="s__item">')
            .append('<a href="/adm/candidates/' + item.Id + '"><span class="name">' + item.Name + " " + item.LName + '</span> <span class="dep">' + item.Dep + "</span></a>")
            .appendTo(ul);
    };

    $('.navbar .navbar-form').submit(function() {
        return false;
    });


    $('.lang-upload').each(function() {
        var id = $(this).attr('data-id');
        $(this).dropzone({
            url: '/adm/langs/' + id + '/upload',
            acceptedFiles: ".json",
            previewsContainer: '#langs-previews',
            uploadMultiple: false,
            clickable: true,
            params: {
                '_xsrf': xsrf
            },
            success: function(i, response) {
                if (response == null) {
                    location.reload()
                } else {
                    show_notice(response.message)
                }
            }
        });
    });

    $('#lang-add-form input[type="file"]').change(function() {
        $('.l__filename', $(this).parents('.form-group')).text($(this).prop('files')[0].name);
    });

    $('body').on('click', '.comment .edit', function() {
        var o = $(this).parents('.comment');
        if ($(this).hasClass('active')) {
            $('.text', o).show();
            $('form', o).hide();
            $(this).removeClass('active');
        } else {
            $('.text', o).hide();
            $('form', o).show();
            $('form textarea', o).focus();
            $(this).addClass('active');
        }
    });

    $('body').on('click', '.comment .remove', function() {
        var t = $(this);
        if (t.hasClass('active')) {
            return false;
        }
        t.addClass('active');

        $.ajax({
            url: t.attr('href'),
            success: function(r) {
                t.removeClass('active');
                if (r === null) {
                    location.reload();
                    return false;
                }
                if (r.success) {
                    t.parents('.comment').remove();
                } else {
                    show_notice(r.message);
                }
            },
            error: function() {
                t.removeClass('active');
            }
        });

        return false;
    });

    $('body').on('click', '.comment form .form-group button:nth-child(2)', function() {
        $('.edit', $(this).parents('.comment')).click();
    });


    $('body').on('submit', '.comment form', function() {
        var t = $(this);
        var o = t.parents('.comment');
        if ($("button[type='submit']", t).hasClass("active")) {
            return false;
        }
        $("button[type='submit']", t).addClass("active");

        $.ajax({
            url: t.attr('action'),
            type: 'POST',
            data: t.serialize(),
            success: function(r) {
                $("button[type='submit']", t).removeClass("active");
                if (r === null) {
                    location.reload();
                    return false;
                }
                if (r.success) {
                    $('.has-error label', t).text('');
                    $('.date', o).text(r.dt);
                    $('.text', o).html(r.text).show();
                    $('.edit', o).removeClass('active');
                    t.hide();
                } else {
                    $('.has-error label', t).text(r.error);
                }
            },
            error: function() {
                $("button[type='submit']", t).removeClass("active");
            }
        });

        return false;
    });


    $('#issue_form').find("[name=_xsrf]").val(xsrf);
    $.ajax({
        url: '/adm/issues/labels',
        type: 'GET',
        success: function(data) {
            if (data !== null && data !== undefined && data.success === undefined) {
                var html = "";
                for (i in data) {
                    html += '<label><input type="checkbox" name="labels[]" value="' + data[i].name + '">' + data[i].name + '</label><br/>';
                }
                $('#issue_labels').html(html);
            }
        },
        error: function() {
            console.log("failed");
        }
    });

    $('#issue_form').submit(function(event) {
        $.ajax({
            url: '/adm/issues/report',
            type: 'POST',
            data: $('#issue_form').serialize(),
            success: function(data) {
                show_notice(data.message);
                if (data.success == true) {
                    $("#create_issue_modal").modal('hide');
                    $('#issue_form').trigger('reset');
                }
            },
            error: function() {
                show_notice("Error while sending request.")
            }
        });

        return false;
    });

});

function show_notice(s) {
    $.snackbar({
        content: s,
        style: "snackbar",
        timeout: 5000
    });
}

function getAjaxData(url, cb) {
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
}

function base64_decode(data) {
    //  discuss at: http://phpjs.org/functions/base64_decode/
    // original by: Tyler Akins (http://rumkin.com)
    // improved by: Thunder.m
    // improved by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
    // improved by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
    //    input by: Aman Gupta
    //    input by: Brett Zamir (http://brett-zamir.me)
    // bugfixed by: Onno Marsman
    // bugfixed by: Pellentesque Malesuada
    // bugfixed by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
    //   example 1: base64_decode('S2V2aW4gdmFuIFpvbm5ldmVsZA==');
    //   returns 1: 'Kevin van Zonneveld'
    //   example 2: base64_decode('YQ===');
    //   returns 2: 'a'

    var b64 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=';
    var o1, o2, o3, h1, h2, h3, h4, bits, i = 0,
        ac = 0,
        dec = '',
        tmp_arr = [];

    if (!data) {
        return data;
    }

    data += '';

    do { // unpack four hexets into three octets using index points in b64
        h1 = b64.indexOf(data.charAt(i++));
        h2 = b64.indexOf(data.charAt(i++));
        h3 = b64.indexOf(data.charAt(i++));
        h4 = b64.indexOf(data.charAt(i++));

        bits = h1 << 18 | h2 << 12 | h3 << 6 | h4;

        o1 = bits >> 16 & 0xff;
        o2 = bits >> 8 & 0xff;
        o3 = bits & 0xff;

        if (h3 == 64) {
            tmp_arr[ac++] = String.fromCharCode(o1);
        } else if (h4 == 64) {
            tmp_arr[ac++] = String.fromCharCode(o1, o2);
        } else {
            tmp_arr[ac++] = String.fromCharCode(o1, o2, o3);
        }
    } while (i < data.length);

    dec = tmp_arr.join('');

    return dec.replace(/\0+$/, '');
}
