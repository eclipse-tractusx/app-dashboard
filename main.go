/********************************************************************************
* Copyright (c) 2023 SAP SE
* Copyright (c) 2023 Contributors to the Eclipse Foundation
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
  ********************************************************************************/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

type templateValues struct {
	Res             Applications
	LastSync        time.Time
	InitialSync     bool
	IgnoreNamespace map[string]bool
	Environment     string
	GitVersion      string
	AppVersion      int
}

var errorPage []byte

const k8sRefreshTimeInSeconds = 300 // 5 Minutes

func main() {
	cluster, envName, ignoreNamespace := initializeAppFlags()

	var err error
	errorPage, err = os.ReadFile("web/error.html")
	if err != nil {
		panic(err)
	}

	clientSet := GetClientSet(*cluster)

	version, err := clientSet.DiscoveryClient.ServerVersion()
	if err != nil {
		log.Fatal(err)
	}

	values := templateValues{
		Res:             Applications{},
		LastSync:        time.Now(),
		InitialSync:     false,
		IgnoreNamespace: ignoreNamespace,
		Environment:     envName,
		GitVersion:      version.GitVersion,
		AppVersion:      1,
	}

	go startWebserver(&values)

	go refreshApplications(&values, k8sRefreshTimeInSeconds, clientSet, ignoreNamespace)

	time.Sleep(time.Duration(1<<63 - 1))
}

func refreshApplications(values *templateValues, refreshTimeInSeconds float64, clientset *kubernetes.Clientset, ignoreNamespace map[string]bool) {
	for true {
		values.Res = requestAndTransformApplications(clientset, ignoreNamespace)

		values.LastSync = time.Now()
		values.InitialSync = true

		time.Sleep(time.Duration(refreshTimeInSeconds * float64(time.Second)))
	}
}

func requestAndTransformApplications(clientset *kubernetes.Clientset, ignoreNamespace map[string]bool) Applications {
	var applicationsResponse = Applications{}

	d, err := clientset.RESTClient().Get().AbsPath("/apis/argoproj.io/v1alpha1/applications").DoRaw(context.TODO())
	if err != nil {
		if statusError, ok := err.(*errors.StatusError); ok {
			if statusError.Status().Code == 404 {
				println("No applications found.")
			} else {
				fmt.Printf("Got an error from the k8s api: %v", string(d))
				panic(statusError)
			}
		} else {
			fmt.Printf("Got an error from the k8s api: %v", string(d))
			panic(err)
		}
	} else {
		if err := json.Unmarshal(d, &applicationsResponse); err != nil {
			panic(err)
		}
		// TODO: Prints debug info on response data; Helpful for seeing what data is available; Should be set to debug
		// fmt.Println(applicationsResponse)
	}

	transformApplicationsResponse(applicationsResponse, ignoreNamespace)

	return applicationsResponse
}

func transformApplicationsResponse(applicationsResponse Applications, ignoreNamespace map[string]bool) {
	for i, item := range applicationsResponse.Items {
		applicationsResponse.Items[i].IgnoreNamespace = false
		if _, ok := ignoreNamespace[applicationsResponse.Items[i].Spec.Destination.Namespace]; ok {
			applicationsResponse.Items[i].IgnoreNamespace = true
		}

		applicationsResponse.Items[i].Status.Summary.LatestImage = false
		applicationsResponse.Items[i].Status.Summary.PostgresqlImageFound = false
		for _, image := range item.Status.Summary.Images {
			if strings.Contains(image, ":latest") || strings.Contains(image, ":main") {
				applicationsResponse.Items[i].Status.Summary.LatestImage = true
			}

			if strings.Contains(strings.ToLower(image), "/postgresql:") {
				parts := strings.Split(image, ":")

				if len(parts) == 2 {
					applicationsResponse.Items[i].Status.Summary.PostgresqlImage = parts[1]
				}
			}
		}
	}
}

func initializeAppFlags() (*bool, string, map[string]bool) {
	cluster := flag.Bool("in-cluster", false, "Specify if the code is running inside a cluster or from outside.")
	flag.Parse()

	envName := "Unset"
	envNameFromENV := strings.TrimSpace(os.Getenv("ENVIRONMENT_NAME"))

	if len(envNameFromENV) > 0 {
		envName = envNameFromENV
	}

	ignoreNamespaceRaw := os.Getenv("IGNORE_NAMESPACE")
	var ignoreNamespace = map[string]bool{}

	if len(ignoreNamespaceRaw) > 0 {
		split := strings.Split(ignoreNamespaceRaw, ",")
		for _, namespace := range split {
			ignoreNamespace[strings.TrimSpace(namespace)] = true
		}
	}

	return cluster, envName, ignoreNamespace
}

func startWebserver(values *templateValues) {
	templates := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"argoHealth": argoHealthToHtmlFunc(),
		"argoSync":   argoSyncStatusToHtmlFunc(),
		"fixGithubUrl": func(url string) string {
			return strings.TrimSuffix(strings.ReplaceAll(url, "git@github.com:", "https://github.com/"), ".git")
		},
		"lastAppSyncShort": func(history []history) string {
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
		"lastAppSyncLong": lastAppSyncToHtmlFunc(),
		"lastSync": func(lastUpdate time.Time) string {

			duration := time.Now().Sub(lastUpdate).Round(time.Second)

			return fmt.Sprint(duration)
		},
		"ignoreNamespace": func(ignoreNamespace map[string]bool) string {

			if len(ignoreNamespace) == 0 {
				return "No Namespaces are set to be ignored."
			}

			result := fmt.Sprintf("(%d): ", len(ignoreNamespace))
			for k, _ := range ignoreNamespace {
				result += k

				result += ", "
			}

			return strings.TrimSuffix(result, ", ")
		},
		"image": containerImageToHtmlFunc(),
	}).ParseFiles("./web/template/index.html"))

	http.Handle("/css/", maxAgeHandler(86400, http.StripPrefix("/css/",
		http.FileServer(http.Dir("./web/css")))))

	http.Handle("/img/", maxAgeHandler(86400, http.StripPrefix("/img/",
		http.FileServer(http.Dir("./web/img")))))

	http.Handle("/webfonts/", maxAgeHandler(86400, http.StripPrefix("/webfonts/",
		http.FileServer(http.Dir("./web/webfonts")))))

	http.Handle("/js/", maxAgeHandler(86400, http.StripPrefix("/js/",
		http.FileServer(http.Dir("./web/js")))))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")

		if values.InitialSync {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI != "/" && r.RequestURI != "/index" && r.RequestURI != "/index.html" {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "text/html")
			w.Write(errorPage)

			return
		}

		w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")

		w.WriteHeader(http.StatusOK)

		if err := templates.ExecuteTemplate(w, "index.html", values); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func maxAgeHandler(seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}

func GetClientSet(cluster bool) *kubernetes.Clientset {
	if cluster == false {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return clientset
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
