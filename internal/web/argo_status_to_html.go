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

const (
	defaultArgoHealthTemplate     = `<i title="Error" class="fa fa-question-circle" style="color: rgb(233, 109, 118);"></i>`
	defaultArgoSyncStatusTemplate = `<i title="Error" class="fa fa-question-circle" style="color: rgb(233, 109, 118);"></i>`
)

var (
	argoHealthToHtmlTemplate = map[string]string{
		"Healthy":     `<i title="Healthy" class="fa-solid fa-heart" style="color: rgb(24, 190, 148);"></i>`,
		"Progressing": `<i title="Progressing" class="fa fa fa-circle-notch" style="color: rgb(13, 173, 234);"></i>`,
		"Degraded":    `<i title="Degraded" class="fa fa-heart-broken" style="color: rgb(233, 109, 118);"></i>`,
		"Suspended":   `<i title="Suspended" class="fa fa-pause-circle" style="color: rgb(118, 111, 148);"></i>`,
		"Missing":     `<i title="Missing" class="fa fa-ghost" style="color: rgb(244, 192, 48);"></i>`,
		"Unknown":     `<i title="Unknown" class="fa fa-question-circle" style="color: rgb(204, 214, 221);"></i>`,
	}
	argoSyncStatusToHtmlTemplate = map[string]string{
		"Synced":    `<i title="Synced" class="fa fa-check-circle" style="color: rgb(24, 190, 148);"></i>`,
		"OutOfSync": `<i title="OutOfSync" class="fa fa-arrow-alt-circle-up" style="color: rgb(244, 192, 48);"></i>`,
	}
)

func ArgoHealthToHtmlFunc() func(status string) string {
	return func(status string) string {
		htmlTemplate, found := argoHealthToHtmlTemplate[status]

		if found {
			return htmlTemplate
		} else {
			return defaultArgoHealthTemplate
		}
	}
}

func ArgoSyncStatusToHtmlFunc() func(status string) string {
	return func(status string) string {
		htmlTemplate, found := argoSyncStatusToHtmlTemplate[status]

		if found {
			return htmlTemplate
		} else {
			return defaultArgoSyncStatusTemplate
		}
	}
}
