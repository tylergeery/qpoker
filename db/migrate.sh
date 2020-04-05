#!/bin/bash

# migrate.sh - Create migration or apply, revert, seed DB

_usage() { 
    echo "./migrate.sh"
    echo "\tPG_CONNECTION=<conn> MSG=<msg> ./migrate.sh migrate"
    echo "\tPG_CONNECTION=<conn> ./migrate.sh apply"
    echo "\tPG_CONNECTION=<conn> ./migrate.sh revert"
    echo "\tPG_CONNECTION=<conn> ./migrate.sh seed"
}

_migrate() {
    echo "Creating migration $GOPATH"
    /go/bin/migrate -path=/go/db/migrations/ -database "$PG_CONNECTION" create -ext sql -dir /go/db/migrations "$MSG"
}

_apply() {
    echo "Applying migrations"
    migrate -path=/go/db/migrations/ -database "$PG_CONNECTION" up
}

_revert() {
    echo "Revert migration"
    migrate -path=/go/db/migrations/ -database "$PG_CONNECTION" down 1
}

_seed() {
    echo "Applying seed data"
    go run seed.go
}

if [ "$1" == "" ]; then
    usage
    exit 1
fi

command="$1"
shift


if [[ "$command" == "migrate" ]]
then
    _migrate
elif [[ "$command" = "apply" ]]
then
    _apply
elif [[ "$command" = "revert" ]]
then
    _revert
elif [[ "$command" = "seed" ]]
then
    _seed
else
    _usage
    exit 1
fi
