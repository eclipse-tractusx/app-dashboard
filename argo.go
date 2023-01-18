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

type history struct {
	DeployStartedAt string
	DeployedAt      string
	Id              int
	Revision        string
	Source          source
}

type status struct {
	Health  health
	History []history
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
