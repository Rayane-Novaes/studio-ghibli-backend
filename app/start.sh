#!/bin/bash
docker run --rm --publish "5555:5432" --name some-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=ghibliApi -v pgdata:/var/lib/postgresql/data postgres 