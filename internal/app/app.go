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
	"time"
)

type Dashboard struct {
	config                   *ApplicationConfig
	web                      Webserver
	gateway                  ApplicationGateway
	syncResult               *ApplicationsSyncResult
	refreshIntervalInSeconds float64
}

func NewDashboard(gateway ApplicationGateway, web Webserver, config *ApplicationConfig) *Dashboard {
	return &Dashboard{
		refreshIntervalInSeconds: 5 * 60,
		gateway:                  gateway,
		web:                      web,
		config:                   config,
		syncResult: &ApplicationsSyncResult{
			Res:             Applications{},
			LastSync:        time.Now(),
			InitialSync:     false,
			IgnoreNamespace: ignoredNamespacesAsMap(config.IgnoredNamespaces),
			Environment:     config.EnvironmentName,
			GitVersion:      "",
			AppVersion:      1,
		},
	}
}

func (d *Dashboard) Run() {
	go d.web.Start(8080, d.syncResult)
	go d.syncApplications()
}

func (d *Dashboard) syncApplications() {
	for true {
		d.syncResult.Res = d.gateway.GetApplications()
		d.syncResult.LastSync = time.Now()
		d.syncResult.InitialSync = true

		time.Sleep(time.Duration(d.refreshIntervalInSeconds * float64(time.Second)))
	}
}

func ignoredNamespacesAsMap(namespaces []string) map[string]bool {
	result := make(map[string]bool, len(namespaces))
	for _, namespace := range namespaces {
		result[namespace] = true
	}
	return result
}
