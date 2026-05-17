FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /registry ./cmd/registry

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /registry /usr/local/bin/registry
WORKDIR /app
EXPOSE 9876
ENTRYPOINT ["/usr/local/bin/registry"]
CMD ["-catalog", "catalog", "-listen", "0.0.0.0:9876"]
