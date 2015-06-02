package ctrls

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"

	"io/ioutil"
	"net/http"
)

type IssueReq struct {
	PID    string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"description"`
	AID    int    `json:"assignee_id"`
	MID    int    `json:"milestone_id"`
	Labels string `json:"labels"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type IssueForm struct {
	Title string `form:"title"`
	Desc  string `form:"description"`
}

type Issue struct {
	BaseController
}

func (c *Issue) GetLabels() {

	authToken := c.GetSession("gitlabToken")

	if authToken == nil {
		c.CheckErr(notNilErr, gitlabAuthFailed)
		return
	}

	if !c.IsAjax() {
		c.Redirect(beego.UrlFor("Cand.Index"), 302)
		return
	}

	res, err := http.Get(gitlabUrl + "projects/" + gitlabProject + "/labels?private_token=" + authToken.(string))
	if c.CheckErr(err, gitlabAuthFailed) {
		return
	}

	if res.Status != "200 OK" {
		c.CheckErr(err, gitlabAuthFailed, "Get Labels error, expected status: 200 OK, got: "+res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if c.CheckErr(err, internalErr) {
		return
	}

	labels := []Label{}

	if c.CheckErr(json.Unmarshal(body, &labels), marshalErr) {
		return
	}
	c.Data["json"] = labels
}

func (c *Issue) ReportIssue() {
	if !c.IsAjax() {
		c.Redirect(beego.UrlFor("Cand.Index"), 302)
		return
	}

	var fLabels string

	form := IssueForm{}
	if c.CheckErr(c.ParseForm(&form), internalErr) {
		return
	}
	labels := c.GetStrings("labels[]")
	fLabels = strings.Join(labels, ",")

	authToken := c.GetSession("gitlabToken")

	if authToken == nil {
		c.CheckErr(notNilErr, gitlabTokenReq, "Auth token for gitlab required")
		return
	}

	issueReq := IssueReq{}
	issueReq.PID = gitlabProject
	issueReq.Title = form.Title
	issueReq.Desc = form.Desc
	issueReq.Labels = fLabels

	data, err := json.Marshal(&issueReq)
	if c.CheckErr(err, marshalErr) {
		return
	}

	res, err := http.Post(gitlabUrl+"projects/"+gitlabProject+"/issues?private_token="+authToken.(string), "application/json", bytes.NewBuffer(data))
	if c.CheckErr(err, gitlabAuthFailed) {
		return
	}

	if res.Status != "201 Created" {
		c.CheckErr(notNilErr, gitlabAuthFailed, res.Status)
		return
	}

	c.Data["json"] = RJson{T("issue_create_success"), true}
}
