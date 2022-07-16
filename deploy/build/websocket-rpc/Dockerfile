FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH="amd64" \
    GOOS=linux

# Move to working directory /build
WORKDIR /go/src/github.com/IM-Lite
RUN git clone https://github.91chi.fun//https://github.com/IM-Lite/IM-Lite-Server.git
WORKDIR /go/src/github.com/IM-Lite/IM-Lite-Server
RUN go mod download

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install

# Build the application
WORKDIR /go/src/github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket
RUN go build --ldflags "-extldflags -static" -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /go/src/github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/main .

# Build a small image
FROM scratch
WORKDIR /app
COPY --from=builder /dist/main /app/websocket
RUN chmod +x /app/websocket && mkdir /app/etc
VOLUME /app/etc/websocket.yaml
CMD ["/app/websocket"]
