/*******************************************************************************
 * Copyright (c) 2021,2023 Contributors to the Eclipse Foundation
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Apache License, Version 2.0 which is available at
 * https://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 ******************************************************************************/

package web

import (
	"dashboard/internal/app"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

type Webserver struct {
	errorPage []byte
}

func NewWebserver() *Webserver {
	errorPage, err := os.ReadFile("web/error.html")
	if err != nil {
		panic(err)
	}

	return &Webserver{errorPage: errorPage}
}

func (web *Webserver) Start(port int, syncResult *app.ApplicationsSyncResult) {
	configureStaticContentServe()

	web.configureRootHandler(createHtmlTemplate(), syncResult)

	log.Printf("Listening on port :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func (web *Webserver) configureRootHandler(template *template.Template, syncResult *app.ApplicationsSyncResult) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI != "/" && r.RequestURI != "/index" && r.RequestURI != "/index.html" {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "text/html")
			w.Write(web.errorPage)

			return
		}

		w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")

		w.WriteHeader(http.StatusOK)

		if err := template.ExecuteTemplate(w, "index.html", syncResult); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func createHtmlTemplate() *template.Template {
	templates := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"argoHealth": ArgoHealthToHtmlFunc(),
		"argoSync":   ArgoSyncStatusToHtmlFunc(),
		"fixGithubUrl": func(url string) string {
			return strings.TrimSuffix(strings.ReplaceAll(url, "git@github.com:", "https://github.com/"), ".git")
		},
		"lastAppSyncShort": func(history []app.History) string {
			sort.Slice(history, func(i, j int) bool {
				return history[i].Id > history[j].Id
			})

			if len(history) < 1 {
				return "none"
			}

			t, _ := time.Parse("2006-01-02T15:04:05Z07:00", history[0].DeployedAt)
			duration := time.Now().Sub(t).Round(time.Minute)

			if duration.Hours() > 24 {
				return fmt.Sprintf("%v days", int(duration.Hours()/24))
			}

			return fmt.Sprint(duration)
		},
		"lastAppSyncLong": LastAppSyncToHtmlFunc(),
		"lastSync": func(lastUpdate time.Time) string {

			duration := time.Now().Sub(lastUpdate).Round(time.Second)

			return fmt.Sprint(duration)
		},
		"ignoreNamespace": func(ignoreNamespace map[string]bool) string {

			if len(ignoreNamespace) == 0 {
				return "No Namespaces are set to be ignored."
			}

			result := fmt.Sprintf("(%d): ", len(ignoreNamespace))
			for k := range ignoreNamespace {
				result += k

				result += ", "
			}

			return strings.TrimSuffix(result, ", ")
		},
		"image": ContainerImageToHtmlFunc(),
	}).ParseFiles("./web/template/index.html"))

	return templates
}

func configureStaticContentServe() {
	http.Handle("/css/", maxAgeHandler(86400, http.StripPrefix("/css/",
		http.FileServer(http.Dir("./web/css")))))

	http.Handle("/img/", maxAgeHandler(86400, http.StripPrefix("/img/",
		http.FileServer(http.Dir("./web/img")))))

	http.Handle("/webfonts/", maxAgeHandler(86400, http.StripPrefix("/webfonts/",
		http.FileServer(http.Dir("./web/webfonts")))))

	http.Handle("/js/", maxAgeHandler(86400, http.StripPrefix("/js/",
		http.FileServer(http.Dir("./web/js")))))
}

func maxAgeHandler(seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}
