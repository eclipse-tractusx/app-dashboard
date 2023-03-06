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

package app

import (
	"flag"
	"os"
	"strings"
)

type Dashboard struct {
	RunsInCluster     *bool
	IgnoredNamespaces []string
	EnvironmentName   string
}

func NewDashboard() *Dashboard {
	dashboard := new(Dashboard)

	dashboard.RunsInCluster = flag.Bool("in-cluster", false, "Specify if the code is running inside a cluster or from outside.")
	dashboard.IgnoredNamespaces = getIgnoredNamespaces()
	dashboard.EnvironmentName = getEnvironmentName()

	flag.Parse()
	return dashboard
}

func getEnvironmentName() string {
	envNameFromENV := strings.TrimSpace(os.Getenv("ENVIRONMENT_NAME"))

	if envNameFromENV != "" {
		return envNameFromENV
	}
	return "Unset"
}

func getIgnoredNamespaces() []string {
	ignoreNamespaceRaw := strings.TrimSpace(os.Getenv("IGNORE_NAMESPACE"))

	if ignoreNamespaceRaw != "" {
		return strings.Split(ignoreNamespaceRaw, ",")
	}
	return []string{}
}
