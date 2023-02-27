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

func argoHealthToHtmlFunc() func(status string) string {
	return func(status string) string {

		// Took the code directly from argocd webui + chrome inspect; Due to use of fontawesome, this works
		switch status {
		case "Healthy":
			return "<i title=\"Healthy\" class=\"fa-solid fa-heart\" style=\"color: rgb(24, 190, 148);\"></i>"
		case "Progressing":
			return "<i title=\"Progressing\" class=\"fa fa fa-circle-notch\" style=\"color: rgb(13, 173, 234);\"></i>"
		case "Degraded":
			return "<i title=\"Degraded\" class=\"fa fa-heart-broken\" style=\"color: rgb(233, 109, 118);\"></i>"
		case "Suspended":
			return "<i title=\"Suspended\" class=\"fa fa-pause-circle\" style=\"color: rgb(118, 111, 148);\"></i>"
		case "Missing":
			return "<i title=\"Missing\" class=\"fa fa-ghost\" style=\"color: rgb(244, 192, 48);\"></i>"
		case "Unknown":
			return "<i title=\"Unknown\" class=\"fa fa-question-circle\" style=\"color: rgb(204, 214, 221);\"></i>"
		default:
			return "<i title=\"Error\" class=\"fa fa-question-circle\" style=\"color: rgb(233, 109, 118);\"></i>"
		}
	}
}
