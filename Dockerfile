FROM golang:1.10 AS build
WORKDIR /go/src
COPY . .
ENV CGO_ENABLED=0
RUN go get -u github.com/gorilla/mux && go build -a --installsuffix cgo --ldflags="-s" -o /go/bin/app

FROM alpine
WORKDIR /app
COPY --from=build /go/bin/app .
ENTRYPOINT ./app