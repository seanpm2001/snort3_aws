# Build the manager binary
FROM golang:1.18 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY ["src/manager/go.mod", "src/manager/go.sum", "src/manager/"]
COPY ["src/message/go.mod", "src/message/go.sum", "src/message/"]
COPY ["src/ipspolicy/go.mod", "src/ipspolicy/go.sum", "src/ipspolicy/"]

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
WORKDIR /workspace/src/manager
RUN go mod download

# Copy the go source
COPY src/message /workspace/src/message
COPY src/ipspolicy /workspace/src/ipspolicy
COPY src/manager/api ./api
COPY src/manager/main.go ./main.go
COPY src/manager/controllers ./controllers

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/src/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
