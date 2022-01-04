#!/bin/sh
set -e

litestream restore -if-db-not-exists -if-replica-exists $APP_DB_CONNECTION_STRING
litestream replicate
