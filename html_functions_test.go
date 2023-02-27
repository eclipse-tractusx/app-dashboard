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

package main

import "testing"

const (
	defaultArgoHealth = `<i title="Error" class="fa fa-question-circle" style="color: rgb(233, 109, 118);"></i>`
)

var (
	argoHealthToHtml = argoHealthToHtmlFunc()
	statusToRender   = ""
	renderedHtml     = ""
)

func TestShouldRenderErrorTemplateForEmptyArgoStatus(t *testing.T) {
	givenStatusToRender("")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(defaultArgoHealth, t)
}

func TestShouldRenderErrorTemplateForInvalidArgoStatus(t *testing.T) {
	givenStatusToRender("blabla")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(defaultArgoHealth, t)
}

func TestShouldRenderHealthyArgoStatus(t *testing.T) {
	givenStatusToRender("Healthy")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Healthy" class="fa-solid fa-heart" style="color: rgb(24, 190, 148);"></i>`, t)
}

func TestShouldRenderProgressingStatus(t *testing.T) {
	givenStatusToRender("Progressing")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Progressing" class="fa fa fa-circle-notch" style="color: rgb(13, 173, 234);"></i>`, t)
}

func TestShouldRenderDegradedStatus(t *testing.T) {
	givenStatusToRender("Degraded")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Degraded" class="fa fa-heart-broken" style="color: rgb(233, 109, 118);"></i>`, t)
}

func TestShouldRenderSuspendedStatus(t *testing.T) {
	givenStatusToRender("Suspended")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Suspended" class="fa fa-pause-circle" style="color: rgb(118, 111, 148);"></i>`, t)
}

func TestShouldRenderMissingStatus(t *testing.T) {
	givenStatusToRender("Missing")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Missing" class="fa fa-ghost" style="color: rgb(244, 192, 48);"></i>`, t)
}

func TestShouldRenderUnknownStatus(t *testing.T) {
	givenStatusToRender("Unknown")

	whenRenderingStatusAsHtml()

	thenRenderedHtmlIs(`<i title="Unknown" class="fa fa-question-circle" style="color: rgb(204, 214, 221);"></i>`, t)
}

func givenStatusToRender(status string) {
	statusToRender = status
}

func whenRenderingStatusAsHtml() {
	renderedHtml = argoHealthToHtml(statusToRender)
}

func thenRenderedHtmlIs(expected string, t *testing.T) {
	if renderedHtml != expected {
		t.Errorf("Argo Status not rendered correctly in HTML!\n Given status: %s\n expected: %s\n actual:%s\n",
			statusToRender, expected, renderedHtml)
	}
}
