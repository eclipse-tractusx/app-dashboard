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

import "strings"

func containerImageToHtmlFunc() func(fullImageUrl string) string {
	return func(fullImageUrl string) string {
		if fullImageUrl == "" {
			return ""
		}

		var tag string
		var image string

		fullUrlSplitByColon := strings.Split(fullImageUrl, ":")

		// i.e. busybox:latest
		if len(fullUrlSplitByColon) == 2 {
			tag = fullUrlSplitByColon[1]
		}
		image = fullUrlSplitByColon[0]

		var path []string
		imageSplitBySlash := strings.Split(image, "/")

		// i.e. tractusx/app-dashboard or tractusx/traceability/irs/item-relationship-service -> irs/item-relationship-service
		if len(imageSplitBySlash) >= 3 {
			path = imageSplitBySlash[:len(imageSplitBySlash)-2]
			image = imageSplitBySlash[len(imageSplitBySlash)-2] + "/" + imageSplitBySlash[len(imageSplitBySlash)-1]
		}

		host := ""
		if len(path) > 0 && isDomain(path[0]) {
			host = path[0]
			path = path[1:]
		}

		return hostAsHtml(host) + pathAsHtml(path) + imageAsHtml(image) + tagAsHtml(tag)
	}
}

func tagAsHtml(tag string) string {
	return `:<span class="tag">` + tag + `</span>`
}

func imageAsHtml(image string) string {
	return `<span class="image">` + image + `</span>`
}

func hostAsHtml(host string) string {
	if host == "" {
		return ""
	}
	return `<span class="host">` + host + `</span>/`
}

func pathAsHtml(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return `<span class="path">` + strings.Join(path, "/") + `</span>/`
}

func isDomain(value string) bool {
	return strings.Contains(value, ".")
}
