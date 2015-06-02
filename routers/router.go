package routers

import (
	"hrkb/ctrls"
	"hrkb/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "github.com/lib/pq"
)

const mustAuthorize = true

var FilterAdmin = func(ctx *context.Context) bool {
	role, _ := ctx.Input.Session("role").(int)
	return role == models.RoleAdmin
}

var FilterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("uid").(int)
	if !ok {
		ctx.Redirect(302, "/login")
	}
}

func init() {

	beego.Router("/", &ctrls.Cand{}, "get:Index")
	beego.InsertFilter("/", beego.BeforeRouter, FilterUser)

	beego.Router("/login", &ctrls.Main{})
	beego.Router("/logout", &ctrls.Main{}, "get:Logout")

	ns :=
		beego.NewNamespace("/adm",
			beego.NSBefore(func(ctx *context.Context) {

				_, ok := ctx.Input.Session("uid").(int)

				if !ok {
					ctx.Redirect(302, "/login")
				}

			}),
			beego.NSNamespace("/users",
				beego.NSCond(FilterAdmin),
				beego.NSRouter("/list", &ctrls.User{}, "get:Index"),
				beego.NSRouter("/add", &ctrls.User{}, "get:Add;post:Post"),
				beego.NSRouter("/:id:int", &ctrls.User{}),
				beego.NSRouter("/:id:int/remove", &ctrls.User{}, "get:Remove"),
			),
			beego.NSNamespace("/dep",
				beego.NSRouter("/list", &ctrls.Dep{}, "get:Index"),
				beego.NSRouter("/add", &ctrls.Dep{}, "get:Add;post:Post"),
				beego.NSRouter("/:id:int", &ctrls.Dep{}),
				beego.NSRouter("/:id:int/remove", &ctrls.Dep{}, "get:Remove"),
			),
			beego.NSNamespace("/crit",
				beego.NSRouter("/list", &ctrls.Crit{}, "get:Index"),
				beego.NSRouter("/:id:int", &ctrls.Crit{}),
				beego.NSRouter("/add", &ctrls.Crit{}, "get:Add;post:Post"),
				beego.NSRouter("/:id:int/remove", &ctrls.Crit{}, "get:Remove"),
			),
			beego.NSNamespace("/candidates",
				beego.NSRouter("/list", &ctrls.Cand{}, "get:Index"),
				beego.NSRouter("/add", &ctrls.Cand{}, "get:Add;post:Post"),
				beego.NSRouter("/:id:int", &ctrls.Cand{}),
				beego.NSRouter("/:id:int/files", &ctrls.Download{}, "get:Index"),
				beego.NSRouter("/:id:int/contacts", &ctrls.Contact{}, "get:Index"),
				beego.NSRouter("/:id:int/contacts/add", &ctrls.Contact{}, "post:Add"),
				beego.NSRouter("/:id:int/contacts/:cid/remove", &ctrls.Contact{}, "get:Remove"),
				beego.NSRouter("/:id:int/edit", &ctrls.Cand{}, "get:Edit"),
				beego.NSRouter("/:id:int/remove", &ctrls.Cand{}, "get:Remove"),
			),
			beego.NSNamespace("/ratings",
				beego.NSRouter("/:id:int", &ctrls.Rat{}),
				beego.NSRouter("/detailed/:id:int", &ctrls.Rat{}, "get:GetDetailed"),
				beego.NSRouter("/my/:id:int", &ctrls.Rat{}, "get:GetUserRatings"),
				beego.NSRouter("/:id:int/remove", &ctrls.Rat{}, "get:Remove"),
			),
			beego.NSNamespace("/trash",
				beego.NSCond(FilterAdmin),
				beego.NSRouter("/list", &ctrls.Trash{}),
				beego.NSRouter("/restore/:type/:id", &ctrls.Trash{}, "get:Restore"),
			),

			beego.NSNamespace("/langs",
				beego.NSRouter("/list", &ctrls.Lang{}, "get:Index"),
				beego.NSRouter("/:id:int/set_default", &ctrls.Lang{}, "get:Default"),
				beego.NSRouter("/:id:int/remove", &ctrls.Lang{}, "get:Remove"),
				beego.NSRouter("/:id:int/upload", &ctrls.Lang{}, "post:Upload"),
				beego.NSRouter("/:id:int/download", &ctrls.Lang{}, "get:Download"),
				beego.NSRouter("/add", &ctrls.Lang{}, "post:Add"),
			),
			beego.NSNamespace("/issues",
				beego.NSRouter("/labels", &ctrls.Issue{}, "get:GetLabels"),
				beego.NSRouter("/report", &ctrls.Issue{}, "post:ReportIssue"),
			),

			beego.NSNamespace("/comments",
				beego.NSRouter("/:id:int", &ctrls.Comments{}),
				beego.NSRouter("/:id:int/edit", &ctrls.Comments{}, "post:Edit"),
				beego.NSRouter("/:id:int/remove", &ctrls.Comments{}, "get:Remove"),
			),
			beego.NSRouter("/search", &ctrls.Cand{}, "get:Search"),
			beego.NSRouter("/profile", &ctrls.Prof{}),

			beego.NSRouter("/upload/:id:int", &ctrls.Upload{}, "post:Index"),
			beego.NSRouter("/upload/img/:id:int", &ctrls.Upload{}, "post:Img"),
			beego.NSRouter("/download/:id:int", &ctrls.Download{}, "get:Get"),
			beego.NSRouter("/download/:id:int/remove", &ctrls.Download{}, "post:Remove"),
			beego.NSRouter("/ratings/:id:int", &ctrls.Rat{}),
			beego.NSRouter("/ratings/detailed/:id:int", &ctrls.Rat{}, "get:GetDetailed"),
			beego.NSRouter("/ratings/my/:id:int", &ctrls.Rat{}, "get:GetUserRatings"),
			beego.NSRouter("/ratings/:id:int/remove", &ctrls.Rat{}, "get:Remove"),
			beego.NSNamespace("/debug",
				beego.NSRouter("/lastlog", &ctrls.Debug{}, "get:LastLog"),
			),
		)

	beego.AddNamespace(ns)
}
