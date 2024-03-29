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
	"testing"
	"time"
)

func TestShouldRenderNoneForEmptySyncHistory(t *testing.T) {
	expectedResult := "none"

	renderedHtml := lastAppSyncToHtmlFunc()(nil)

	if renderedHtml != expectedResult {
		t.Errorf("Did not render corretly for empty sync history! \nexpected: %s \ngot: %s", expectedResult, renderedHtml)
	}

	renderedHtml = lastAppSyncToHtmlFunc()([]app.History{})

	if renderedHtml != expectedResult {
		t.Errorf("Did not render corretly for empty sync history! \nexpected: %s \ngot: %s", expectedResult, renderedHtml)
	}
}

func TestShouldRenderSyncHistory(t *testing.T) {
	// overwrite currentTime to make rendered HTML results assertable
	currentTime = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2022-09-18T08:00:00.20Z")
		return t
	}
	historyEntry := app.History{
		DeployStartedAt: "2022-09-18T07:25:40.20Z",
		DeployedAt:      "2022-09-18T07:26:00.20Z",
		Id:              1,
		Revision:        "b8d56b2d875b183f3109f645443373e18f56783b",
		Source: app.Source{
			RepoUrl:        "https://github.com/eclipse-tractusx/app-dashboard",
			Path:           "",
			TargetRevision: "3d7377d0af2683eb89f7c572d7f01fa794260e55",
		},
	}

	historyEntries := []app.History{
		historyEntry,
	}
	expectedHtml := `<li>` + historyEntry.DeployedAt + ` (34m0s)<br/>rev: <a href="https://github.com/eclipse-tractusx/app-dashboard/tree/3d7377d0af2683eb89f7c572d7f01fa794260e55">3d7377d0af2683eb89f7c572d7f01fa794260e55</a></li>`

	renderedHtml := lastAppSyncToHtmlFunc()(historyEntries)

	if renderedHtml != expectedHtml {
		t.Errorf("Sync history Entry not rendered correctly! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}

func TestShouldOrderBySyncHistoryId(t *testing.T) {
	currentTime = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2022-09-18T08:00:00.20Z")
		return t
	}
	firstHistoryEntry := app.History{
		DeployStartedAt: "2022-09-18T07:25:40.20Z",
		DeployedAt:      "2022-09-18T07:26:00.20Z",
		Id:              1,
		Revision:        "c4232944d8e75ab1e23067f1cf4c88f51f82317e",
		Source: app.Source{
			RepoUrl:        "https://github.com/eclipse-tractusx/app-dashboard",
			Path:           "",
			TargetRevision: "c4232944d8e75ab1e23067f1cf4c88f51f82317e",
		},
	}
	secondHistoryEntry := app.History{
		DeployStartedAt: "2022-09-18T07:25:40.20Z",
		DeployedAt:      "2022-09-18T07:26:00.20Z",
		Id:              2,
		Revision:        "b8d56b2d875b183f3109f645443373e18f56783b",
		Source: app.Source{
			RepoUrl:        "https://github.com/eclipse-tractusx/app-dashboard",
			Path:           "",
			TargetRevision: "b8d56b2d875b183f3109f645443373e18f56783b",
		},
	}
	thirdHistoryEntry := app.History{
		DeployStartedAt: "2022-09-18T07:25:40.20Z",
		DeployedAt:      "2022-09-18T07:26:00.20Z",
		Id:              3,
		Revision:        "30a8d5a5e31091a0a450ecde84ac4b0bc3f57cef",
		Source: app.Source{
			RepoUrl:        "https://github.com/eclipse-tractusx/app-dashboard",
			Path:           "",
			TargetRevision: "30a8d5a5e31091a0a450ecde84ac4b0bc3f57cef",
		},
	}

	historyEntries := []app.History{
		secondHistoryEntry,
		thirdHistoryEntry,
		firstHistoryEntry,
	}

	expectedHtml := `<li>` + thirdHistoryEntry.DeployedAt + ` (34m0s)<br/>rev: <a href="https://github.com/eclipse-tractusx/app-dashboard/tree/` + thirdHistoryEntry.Source.TargetRevision + `">` + thirdHistoryEntry.Source.TargetRevision + `</a></li>`
	expectedHtml += `<li>` + secondHistoryEntry.DeployedAt + ` (34m0s)<br/>rev: <a href="https://github.com/eclipse-tractusx/app-dashboard/tree/` + secondHistoryEntry.Source.TargetRevision + `">` + secondHistoryEntry.Source.TargetRevision + `</a></li>`
	expectedHtml += `<li>` + firstHistoryEntry.DeployedAt + ` (34m0s)<br/>rev: <a href="https://github.com/eclipse-tractusx/app-dashboard/tree/` + firstHistoryEntry.Source.TargetRevision + `">` + firstHistoryEntry.Source.TargetRevision + `</a></li>`

	renderedHtml := lastAppSyncToHtmlFunc()(historyEntries)

	if renderedHtml != expectedHtml {
		t.Errorf("Sync history not sorted! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}

func TestShouldCutOffGitExtensionsInRepoUrl(t *testing.T) {
	// overwrite currentTime to make rendered HTML results assertable
	currentTime = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2022-09-18T08:00:00.20Z")
		return t
	}
	historyEntry := app.History{
		DeployStartedAt: "2022-09-18T07:25:40.20Z",
		DeployedAt:      "2022-09-18T07:26:00.20Z",
		Id:              1,
		Revision:        "b8d56b2d875b183f3109f645443373e18f56783b",
		Source: app.Source{
			RepoUrl:        "https://github.com/eclipse-tractusx/app-dashboard.git",
			TargetRevision: "3d7377d0af2683eb89f7c572d7f01fa794260e55",
		},
	}

	historyEntries := []app.History{
		historyEntry,
	}
	expectedHtml := `<li>` + historyEntry.DeployedAt + ` (34m0s)<br/>rev: <a href="https://github.com/eclipse-tractusx/app-dashboard/tree/3d7377d0af2683eb89f7c572d7f01fa794260e55">3d7377d0af2683eb89f7c572d7f01fa794260e55</a></li>`

	renderedHtml := lastAppSyncToHtmlFunc()(historyEntries)

	if renderedHtml != expectedHtml {
		t.Errorf("Sync history Entry not rendered correctly! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}

func TestShouldRenderHttpUrlInsteadOfSshUrl(t *testing.T) {
	// overwrite currentTime to make rendered HTML results assertable
	currentTime = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2022-09-18T08:00:00.20Z")
		return t
	}
	historyEntry := app.History{
		DeployStartedAt: "2022-09-18T07:25:40.20Z",
		DeployedAt:      "2022-09-18T07:26:00.20Z",
		Id:              1,
		Revision:        "b8d56b2d875b183f3109f645443373e18f56783b",
		Source: app.Source{
			RepoUrl:        "git@github.com:eclipse-tractusx/app-dashboard",
			TargetRevision: "3d7377d0af2683eb89f7c572d7f01fa794260e55",
		},
	}

	historyEntries := []app.History{
		historyEntry,
	}
	expectedHtml := `<li>` + historyEntry.DeployedAt + ` (34m0s)<br/>rev: <a href="https://github.com/eclipse-tractusx/app-dashboard/tree/3d7377d0af2683eb89f7c572d7f01fa794260e55">3d7377d0af2683eb89f7c572d7f01fa794260e55</a></li>`

	renderedHtml := lastAppSyncToHtmlFunc()(historyEntries)

	if renderedHtml != expectedHtml {
		t.Errorf("Sync history Entry not rendered correctly! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}

func TestShouldNotRenderLinksForSyncsToCentralChartRepo(t *testing.T) {
	// overwrite currentTime to make rendered HTML results assertable
	currentTime = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2022-09-18T08:00:00.20Z")
		return t
	}
	historyEntry := app.History{
		DeployStartedAt: "2022-09-18T07:25:40.20Z",
		DeployedAt:      "2022-09-18T07:26:00.20Z",
		Id:              1,
		Revision:        "3.0.5",
		Source: app.Source{
			RepoUrl:        "https://eclipse-tractusx.github.io/charts/dev",
			TargetRevision: "3.0.5",
		},
	}

	historyEntries := []app.History{
		historyEntry,
	}
	expectedHtml := `<li>` + historyEntry.DeployedAt + ` (34m0s)<br/>rev: 3.0.5</li>`

	renderedHtml := lastAppSyncToHtmlFunc()(historyEntries)

	if renderedHtml != expectedHtml {
		t.Errorf("Sync history Entry not rendered correctly! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}
