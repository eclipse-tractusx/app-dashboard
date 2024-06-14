###############################################################
# Copyright (c) 2023 Contributors to the Eclipse Foundation
#
# See the NOTICE file(s) distributed with this work for additional
# information regarding copyright ownership.
#
# This program and the accompanying materials are made available under the
# terms of the Apache License, Version 2.0 which is available at
# https://www.apache.org/licenses/LICENSE-2.0.
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.
#
# SPDX-License-Identifier: Apache-2.0
###############################################################

FROM golang:1.21.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go test ./... && \
    CGO_ENABLED=0 \
    go build -installsuffix 'static' -ldflags="-w -s" .


FROM alpine:3.18.4 AS final

WORKDIR /app

COPY ./web /app/web
COPY --from=builder --chown=nonroot:nonroot /app/dashboard /app/dashboard

RUN adduser -u 1000 --disabled-password --gecos "" --no-create-home nonroot
USER nonroot

ENTRYPOINT ["/app/dashboard"]

CMD ["-in-cluster=true"]
