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

WORKDIR /go/db

COPY ./db/migrate.sh migrate.sh



FROM base as prod

# TODO: prod