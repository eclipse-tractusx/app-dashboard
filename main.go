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
	"strings"
	"text/template"
	"time"
)

type metadata struct {
	Generation int
	Name       string
	Namespace  string
}

type summary struct {
	ExternalUrls         []string
	Images               []string
	LatestImage          bool
	PostgresqlImageFound bool
	PostgresqlImage      string
}

type health struct {
	Status string
}

type status struct {
	Health  health
	history []string
	Summary summary
	Sync    statusSync
}

type source struct {
	RepoUrl        string
	Path           string
	TargetRevision string
}

type destination struct {
	Namespace string
	Server    string
}

type spec struct {
	Destination destination
	Project     string
	Source      source
}

type statusSync struct {
	Source source
	Status string
}

type item struct {
	ApiVersion      string
	Kind            string
	Metadata        metadata
	Spec            spec
	Status          status
	IgnoreNamespace bool
}

type Applications struct {
	ApiVersion string
	Items      []item
	Kind       string
}

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
			if strings.Contains(image, ":latest") {
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
		"argoHealth": func(status string) string {

			// Took the code directly from argocd webui + chrome inspect; Due to use of fontawesome, this works
			switch status {
			case "Healthy":
				return "<i title=\"Healthy\" class=\"fa-solid fa-heart\" style=\"color: rgb(24, 190, 148);\"></i>"
			case "Progressing":
				return "<i title=\"Progressing\" class=\"fa fa fa-circle-notch\" style=\"color: rgb(13, 173, 234);\"></i>"
			case "Degraded":
				return "<i title=\"Degraded\" class=\"fa fa-heart-broken\" style=\"color: rgb(233, 109, 118);\"></i>"
			case "Suspended":
				return "<i title=\"Suspended\" class=\"fa fa-pause-circle\" style=\"color: rgb(118, 111, 148);\"></i>"
			case "Missing":
				return "<i title=\"Missing\" class=\"fa fa-ghost\" style=\"color: rgb(244, 192, 48);\"></i>"
			case "Unknown":
				return "<i title=\"Unknown\" class=\"fa fa-question-circle\" style=\"color: rgb(204, 214, 221);\"></i>"
			default:
				return "<i title=\"Error\" class=\"fa fa-question-circle\" style=\"color: rgb(233, 109, 118);\"></i>"
			}
		},
		"argoSync": func(status string) string {

			// Took the code directly from argocd webui + chrome inspect; Due to use of fontawesome, this works
			switch status {
			case "Synced":
				return "<i title=\"Synced\" class=\"fa fa-check-circle\" style=\"color: rgb(24, 190, 148);\"></i>"
			case "OutOfSync":
				return "<i title=\"OutOfSync\" class=\"fa fa-arrow-alt-circle-up\" style=\"color: rgb(244, 192, 48);\"></i>"
			default:
				return "<i title=\"Error\" class=\"fa fa-question-circle\" style=\"color: rgb(233, 109, 118);\"></i>"
			}
		},
		"fixGithubUrl": func(url string) string {
			return strings.TrimSuffix(strings.ReplaceAll(url, "git@github.com:", "https://github.com/"), ".git")
		},
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
		"image": func(fullImageUrl string) string {

			if len(fullImageUrl) == 0 {
				return fullImageUrl
			}

			parts := strings.Split(fullImageUrl, "/")

			if len(parts) == 0 {
				return fullImageUrl
			}

			tags := strings.Split(parts[len(parts)-1], ":")

			var tag string
			var image string
			if len(tags) == 2 {
				tag = tags[1]
				image = tags[0]
			} else {
				return image
			}

			var path string
			for i, part := range parts {
				if i >= len(parts)-2 {
					break
				}

				if i == 0 {
					path += "<span class=\"host\">" + part + "</span>/"
				} else {
					path += "<span class=\"path\">" + part + "</span>/"
				}

			}

			return path + "<span class=\"image\">" + image + "</span>:<span class=\"tag\">" + tag + "</span>"

		},
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
