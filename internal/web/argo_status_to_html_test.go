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

import "testing"

const (
	defaultArgoHealth = `<i title="Error" class="fa fa-question-circle" style="color: rgb(233, 109, 118);"></i>`
)

var (
	healthStatusToRender = ""
	syncStatusToRender   = ""
	renderedHtml         = ""
)

func TestShouldRenderErrorTemplateForEmptyArgoStatus(t *testing.T) {
	givenHealthStatusToRender("")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(defaultArgoHealth, t)
}

func TestShouldRenderErrorTemplateForInvalidArgoStatus(t *testing.T) {
	givenHealthStatusToRender("blabla")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(defaultArgoHealth, t)
}

func TestShouldRenderHealthyArgoStatus(t *testing.T) {
	givenHealthStatusToRender("Healthy")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Healthy" class="fa-solid fa-heart" style="color: rgb(24, 190, 148);"></i>`, t)
}

func TestShouldRenderProgressingStatus(t *testing.T) {
	givenHealthStatusToRender("Progressing")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Progressing" class="fa fa fa-circle-notch" style="color: rgb(13, 173, 234);"></i>`, t)
}

func TestShouldRenderDegradedStatus(t *testing.T) {
	givenHealthStatusToRender("Degraded")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Degraded" class="fa fa-heart-broken" style="color: rgb(233, 109, 118);"></i>`, t)
}

func TestShouldRenderSuspendedStatus(t *testing.T) {
	givenHealthStatusToRender("Suspended")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Suspended" class="fa fa-pause-circle" style="color: rgb(118, 111, 148);"></i>`, t)
}

func TestShouldRenderMissingStatus(t *testing.T) {
	givenHealthStatusToRender("Missing")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Missing" class="fa fa-ghost" style="color: rgb(244, 192, 48);"></i>`, t)
}

func TestShouldRenderUnknownStatus(t *testing.T) {
	givenHealthStatusToRender("Unknown")

	whenRenderingHealthStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Unknown" class="fa fa-question-circle" style="color: rgb(204, 214, 221);"></i>`, t)
}

func TestShouldRenderErrorTemplateForEmptySyncStatus(t *testing.T) {
	givenSyncStatusToRender("")

	whenRenderingSyncStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Error" class="fa fa-question-circle" style="color: rgb(233, 109, 118);"></i>`, t)
}

func TestShouldRenderSyncedSyncStatus(t *testing.T) {
	givenSyncStatusToRender("Synced")

	whenRenderingSyncStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Synced" class="fa fa-check-circle" style="color: rgb(24, 190, 148);"></i>`, t)
}

func TestShouldRenderOutOfSyncStatus(t *testing.T) {
	givenSyncStatusToRender("OutOfSync")

	whenRenderingSyncStatusAsHtml()

	thenRenderedHtmlIs(`<i title="OutOfSync" class="fa fa-arrow-alt-circle-up" style="color: rgb(244, 192, 48);"></i>`, t)
}

func givenHealthStatusToRender(status string) {
	healthStatusToRender = status
}

func givenSyncStatusToRender(status string) {
	syncStatusToRender = status
}

func whenRenderingSyncStatusAsHtml() {
	renderedHtml = argoSyncStatusToHtmlFunc()(syncStatusToRender)
}

func whenRenderingHealthStatusAsHtml() {
	renderedHtml = argoHealthToHtmlFunc()(healthStatusToRender)
}

func thenRenderedHtmlIs(expected string, t *testing.T) {
	if renderedHtml != expected {
		t.Errorf("Argo Status not rendered correctly in HTML!\n Given status: %s\n expected: %s\n actual:%s\n",
			healthStatusToRender, expected, renderedHtml)
	}
}
