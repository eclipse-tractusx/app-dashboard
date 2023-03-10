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
	"dashboard/internal/app"
	"dashboard/internal/gateway"
	"dashboard/internal/web"
	"os"
	"strings"
	"time"
)

func main() {
	dashboard := app.NewDashboard(gateway.NewApplicationGateway(), web.NewWebserver(), getAppConfig())
	dashboard.Run()

	time.Sleep(time.Duration(1<<63 - 1))
}

func getAppConfig() *app.ApplicationConfig {
	return &app.ApplicationConfig{
		IgnoredNamespaces: getIgnoredNamespaces(),
		EnvironmentName:   getEnvironmentName(),
	}
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
