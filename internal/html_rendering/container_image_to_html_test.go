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

package html_rendering

import "testing"

func TestShouldRenderNothingForEmptyImageUrl(t *testing.T) {
	renderedHtml := ContainerImageToHtmlFunc()("")

	if renderedHtml != "" {
		t.Errorf("Did render something for empty image. Nothing expected! Got: %s", renderedHtml)
	}
}

func TestShouldRenderOfficialImagesUnchanged(t *testing.T) {
	officialImageName := "busybox:latest"
	expectedHtml := `<span class="image">busybox</span>:<span class="tag">latest</span>`

	renderedHtml := ContainerImageToHtmlFunc()(officialImageName)

	if renderedHtml != expectedHtml {
		t.Errorf("Did not render official image correctly! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}

func TestShouldRenderImagesWithNamespace(t *testing.T) {
	imageWithNamespace := "tractusx/app-dashboard:1.0.0"
	expectedHtml := `<span class="image">tractusx/app-dashboard</span>:<span class="tag">1.0.0</span>`

	renderedHtml := ContainerImageToHtmlFunc()(imageWithNamespace)

	if renderedHtml != expectedHtml {
		t.Errorf("Did not render image with namespace! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}

func TestShouldRenderImageFromAnyContainerRegistry(t *testing.T) {
	ghcrImage := "ghcr.io/catenax-ng/semantic-hub:0.1.0-M3"
	expectedHtml := `<span class="host">ghcr.io</span>/<span class="image">catenax-ng/semantic-hub</span>:<span class="tag">0.1.0-M3</span>`

	renderedHtml := ContainerImageToHtmlFunc()(ghcrImage)

	if renderedHtml != expectedHtml {
		t.Errorf("Did not render image from non DockerHub registry! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}

func TestShouldRenderImageWithMultipleNamespaces(t *testing.T) {
	imageWithMultipleNamespaces := "tractusx/traceability/irs/item-relationship-service:1.0.0"
	expectedHtml := `<span class="path">tractusx/traceability</span>/<span class="image">irs/item-relationship-service</span>:<span class="tag">1.0.0</span>`

	renderedHtml := ContainerImageToHtmlFunc()(imageWithMultipleNamespaces)

	if renderedHtml != expectedHtml {
		t.Errorf("Did not render image with multiple namespaces! \nexpected: %s \nGot: %s", expectedHtml, renderedHtml)
	}
}
