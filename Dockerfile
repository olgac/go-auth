FROM golang:1.10.3-alpine3.7 AS build
WORKDIR /go/src
COPY . .
ENV CGO_ENABLED=0
RUN go get -u github.com/gorilla/mux github.com/dgrijalva/jwt-go && go get github.com/garyburd/redigo/redis && go build -a --installsuffix cgo --ldflags="-s" -o /go/bin/app

FROM alpine
WORKDIR /app
COPY --from=build /go/bin/app .
ENTRYPOINT ./app