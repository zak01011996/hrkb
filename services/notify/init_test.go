package notify

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	M "hrkb/models"
)

func newFakeServer(code int, body string) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	cl := &http.Client{Transport: tr}

	return server, cl
}

//Function creates fake api server wich sleeps to emulate slow internet connection
func newSleepyServer(timeout int, code int, body string) (*httptest.Server, *http.Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")

		time.Sleep(time.Duration(timeout+500) * time.Millisecond)

		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	cl := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Millisecond,
	}

	return server, cl
}

func newMail(id int) M.Mail {
	mail := M.Mail{}
	mail.Id = id
	mail.ToMail = "shoptoli@gmail.com"
	mail.ToName = sql.NullString{"Abdullo", true}
	mail.ToType = "to"
	mail.FromMail = "ahidoyatov@gmail.com"
	mail.FromName = sql.NullString{"hrkb", true}
	mail.Subject = sql.NullString{"notification", true}
	mail.Html = sql.NullString{"<h1>long interestin html</h1>", true}
	mail.Text = sql.NullString{"some text", true}
	return mail
}
