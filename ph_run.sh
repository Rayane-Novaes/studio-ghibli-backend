#!/bin/bash

echoerr() { echo "$@" 1>&2; }

# Update docs
echo ">> Updating Swagger docs..."
go run github.com/swaggo/swag/cmd/swag init --outputTypes go --propertyStrategy snakecase -pd
if [ $? -ne 0 ]; then
  echoerr "Swagger generation failed, please check the messages above"
  exit 1
fi

function export_all() {
	set -o allexport; source "${1}"; set +o allexport
}

export_all .env
go run main.go
