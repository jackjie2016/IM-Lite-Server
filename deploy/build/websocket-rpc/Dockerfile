FROM registry.cn-shanghai.aliyuncs.com/showurl-public/librdkafka-go-builder:latest AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH="amd64" \
    GOOS=linux \
    GOPROXY="https://goproxy.cn,direct"
# Move to working directory /build
WORKDIR /go/src/github.com/IM-Lite
RUN git clone https://github.91chi.fun//https://github.com/IM-Lite/IM-Lite-Server.git
WORKDIR /go/src/github.com/IM-Lite/IM-Lite-Server
RUN go mod download

# Build the application
WORKDIR /go/src/github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket
RUN go build -tags musl --ldflags "-extldflags -static" -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /go/src/github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/main .

# Build a small image
FROM scratch
WORKDIR /app/etc
COPY --from=builder /dist/main /app/websocket
VOLUME /app/etc/websocket.yaml
WORKDIR /app
ENTRYPOINT ["/app/websocket"]