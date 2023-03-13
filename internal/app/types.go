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

type ApplicationGateway interface {
	GetApplications() Applications
	ToolInfoAsHtml() string
}

type Webserver interface {
	Start(port int, syncResult *ApplicationsSyncResult)
}

type ApplicationConfig struct {
	IgnoredNamespaces []string
	EnvironmentName   string
}

type ApplicationsSyncResult struct {
	Res             Applications
	LastSync        time.Time
	InitialSync     bool
	IgnoreNamespace map[string]bool
	Environment     string
	GitVersion      string
	AppVersion      int
}

type Applications struct {
	ApiVersion string
	Items      []item
	Kind       string
}

type item struct {
	ApiVersion      string
	Kind            string
	Metadata        metadata
	Spec            spec
	Status          status
	IgnoreNamespace bool
}

type metadata struct {
	Generation int
	Name       string
	Namespace  string
}

type spec struct {
	Destination destination
	Project     string
	Source      source
}

type status struct {
	Health  health
	History []History
	Summary summary
	Sync    statusSync
}

type destination struct {
	Namespace string
	Server    string
}

type source struct {
	RepoUrl        string
	Path           string
	TargetRevision string
}

type health struct {
	Status string
}

type History struct {
	DeployStartedAt string
	DeployedAt      string
	Id              int
	Revision        string
	Source          source
}

type summary struct {
	ExternalUrls         []string
	Images               []string
	LatestImage          bool
	PostgresqlImageFound bool
	PostgresqlImage      string
}

type statusSync struct {
	Source source
	Status string
}
