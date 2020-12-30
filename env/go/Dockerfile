FROM golang:1.14 as base

WORKDIR /src



FROM base as dev

# Add auto-reload entry point for 
RUN go get -u github.com/githubnemo/CompileDaemon


FROM dev as http_dev

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/qpoker_http ./http" -command="/go/bin/qpoker_http"


FROM dev as ws_dev

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/qpoker_ws ./websocket" -command="/go/bin/qpoker_ws"




FROM base as migrator

RUN go get -u github.com/lib/pq
RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

COPY ./db /go/db
WORKDIR /go/db



FROM golang:alpine as base_prod

COPY ./src /src
WORKDIR /src

RUN go get -d -v ./...

RUN go build -o /go/bin/http ./http
RUN go build -o /go/bin/ws ./websocket



FROM base_prod as http_prod

CMD ["/go/bin/http"]


FROM base_prod as ws_prod


CMD ["/go/bin/ws"]