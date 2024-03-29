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
	"sort"
	"strings"
	"time"
)

var currentTime = getCurrentTime

func lastAppSyncToHtmlFunc() func(history []app.History) string {
	return func(history []app.History) string {
		sort.Slice(history, func(i, j int) bool {
			return history[i].Id > history[j].Id
		})

		if len(history) < 1 {
			return "none"
		}

		var result string
		for _, entry := range history {

			t, _ := time.Parse("2006-01-02T15:04:05Z07:00", entry.DeployedAt)
			duration := currentTime().Sub(t).Round(time.Minute)

			var since string
			if duration.Hours() > 24 {
				since = fmt.Sprintf("%v days", int(duration.Hours()/24))
			} else {
				since = fmt.Sprintf("%v", duration)
			}

			result += "<li>" + entry.DeployedAt + " (" + since + ")<br/>rev: " + linkToRevision(entry.Source) + "</li>"
		}

		return result
	}
}

func linkToRevision(source app.Source) string {
	// Ignore deployments of released charts from central repo, since there are no tags present in this repo
	// Information about the origin of the released chart (product repo) not available in current data structure
	if strings.Contains(source.RepoUrl, "eclipse-tractusx.github.io/charts") {
		return source.TargetRevision
	}

	return `<a href="` + ensureHttpGitHubUrl(source.RepoUrl) + `/tree/` + source.TargetRevision + `">` + source.TargetRevision + `</a>`
}

func ensureHttpGitHubUrl(url string) string {
	return strings.TrimSuffix(strings.ReplaceAll(url, "git@github.com:", "https://github.com/"), ".git")
}

func getCurrentTime() time.Time {
	return time.Now()
}
