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

package gateway

import (
	"context"
	"dashboard/internal/app"
	"encoding/json"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)

type ApplicationGateway struct {
	clientset         *kubernetes.Clientset
	ignoredNamespaces map[string]bool
}

func NewApplicationGateway() *ApplicationGateway {
	var clientSet *kubernetes.Clientset
	if runsInCluster() {
		clientSet = getClusterClientSet()
	} else {
		clientSet = getLocalClientSet()
	}

	return &ApplicationGateway{clientset: clientSet, ignoredNamespaces: ignoredNamespacesAsMap()}
}

func (gateway *ApplicationGateway) GetApplications() app.Applications {
	var applicationsResponse = app.Applications{}

	d, err := gateway.clientset.RESTClient().Get().AbsPath("/apis/argoproj.io/v1alpha1/applications").DoRaw(context.TODO())
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

	transformApplicationsResponse(applicationsResponse, gateway.ignoredNamespaces)

	return applicationsResponse
}

func (gateway *ApplicationGateway) ToolInfoAsHtml() string {
	clusterVersion := getClusterVersion(gateway)
	ignoredNamespaces := getIgnoredNamespacesRaw()

	return fmt.Sprintf("<ul><li>GitVersion / K8s cluster: %s</li><li>Ignored Namespaces: %s</li></ul>", clusterVersion, ignoredNamespaces)
}

func transformApplicationsResponse(applications app.Applications, ignoreNamespace map[string]bool) {
	for i, item := range applications.Items {
		applications.Items[i].IgnoreNamespace = false
		if _, ok := ignoreNamespace[applications.Items[i].Spec.Destination.Namespace]; ok {
			applications.Items[i].IgnoreNamespace = true
		}

		applications.Items[i].Status.Summary.LatestImage = false
		applications.Items[i].Status.Summary.PostgresqlImageFound = false
		for _, image := range item.Status.Summary.Images {
			if strings.Contains(image, ":latest") || strings.Contains(image, ":main") {
				applications.Items[i].Status.Summary.LatestImage = true
			}

			if strings.Contains(strings.ToLower(image), "/postgresql:") {
				parts := strings.Split(image, ":")

				if len(parts) == 2 {
					applications.Items[i].Status.Summary.PostgresqlImage = parts[1]
				}
			}
		}
	}
}

func getClusterVersion(gateway *ApplicationGateway) string {
	version, err := gateway.clientset.DiscoveryClient.ServerVersion()
	versionString := version.GitVersion
	if err != nil {
		versionString = "unknown"
	}
	return versionString
}

func getIgnoredNamespacesRaw() string {
	return strings.TrimSpace(os.Getenv("IGNORE_NAMESPACE"))
}

func getClusterClientSet() *kubernetes.Clientset {
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

func getLocalClientSet() *kubernetes.Clientset {
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

func runsInCluster() bool {
	inCluster := flag.Bool("in-cluster", false, "Specify if the code is running inside a cluster or from outside.")
	flag.Parse()
	return *inCluster
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func ignoredNamespacesAsMap() map[string]bool {
	var namespaces []string
	ignoredNamespaces := getIgnoredNamespacesRaw()
	result := make(map[string]bool)
	if ignoredNamespaces != "" {
		namespaces = strings.Split(ignoredNamespaces, ",")
	}

	for _, namespace := range namespaces {
		result[namespace] = true
	}
	return result
}
