FROM golang:1.19.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go test ./... && \
    CGO_ENABLED=0 \
    go build -installsuffix 'static' -ldflags="-w -s" .


FROM gcr.io/distroless/static:nonroot AS final

COPY ./web /app/web
COPY --from=builder --chown=nonroot:nonroot /app/dashboard /app/dashboard

ENTRYPOINT ["/app/dashboard"]

CMD ["-in-cluster=true"]