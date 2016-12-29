package main

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web"
	. "eaciit/wfdemo-git/web/controller"
	"eaciit/wfdemo-git/web/helper"
	"fmt"
	"net/http"
	"strings"

	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

func main() {
	config := new(Configuration)

	config.ID = "port"
	port, err := config.GetPort()
	if err != nil {
		toolkit.Printf("Error get port : %s \n", err.Error())
	}

	ConfigPath = CONFIG_PATH
	prefix := web.AppName

	ServerAddress = toolkit.Sprintf("localhost:%v", toolkit.ToString(port))
	appConf := knot.AppContainerConfig{Address: ServerAddress}
	otherRoutes := map[string]knot.FnContent{
		"/": func(r *knot.WebContext) interface{} {
			sessionid := r.Session("sessionid", "")
			if sessionid == "" {
				http.Redirect(r.Writer, r.Request, fmt.Sprintf("/%s/page/login", prefix), 301)
			} else {
				http.Redirect(r.Writer, r.Request, fmt.Sprintf("/%s/page/dashboard", prefix), 301)
			}

			return true
		},
		"prerequest": func(r *knot.WebContext) interface{} {
			r.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			r.Writer.Header().Set("Pragma", "no-cache")
			r.Writer.Header().Set("Expires", "0")

			rURL := r.Request.URL.String()
			sessionid := r.Session("sessionid", "")
			helper.WC = r

			if rURL != "/"+prefix+"/page/login" && rURL != "/"+prefix+"/login/processlogin" {
				active := acl.IsSessionIDActive(toolkit.ToString(sessionid))

				if !active {
					r.SetSession("sessionid", "")
					urlSplit := strings.Split(rURL, "/")

					if len(urlSplit) > 2 {
						if urlSplit[2] == "page" {
							http.Redirect(r.Writer, r.Request, fmt.Sprintf("/%s/page/login", prefix), 301)
						} else {
							helper.WC = r
						}
					} else {
						http.Redirect(r.Writer, r.Request, fmt.Sprintf("/%s/page/login", prefix), 301)
					}
				}
			}
			return true
		},
		"postrequest": func(r *knot.WebContext) interface{} {
			WriteLog(r.Session("sessionid", ""), "access", r.Request.URL.String())
			return true
		},
	}

	knot.DefaultOutputType = knot.OutputTemplate
	knot.StartContainerWithFn(&appConf, otherRoutes)
}
