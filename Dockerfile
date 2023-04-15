FROM golang:1.20.3

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o gRPC-task .cmd/server/main.go

CMD ["./gRPC-task"]