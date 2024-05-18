#!/bin/sh
set -e

CONFIG_FILE="/app/config.ini"

if [ ! -z "$DB_HOST" ]; then
    sed -i "s/^host.*/host = $DB_HOST/" $CONFIG_FILE
fi

if [ ! -z "$DB_PORT" ]; then
    sed -i "s/^port.*/port = $DB_PORT/" $CONFIG_FILE
fi

if [ ! -z "$DB_USER" ]; then
    sed -i "s/^user.*/user = $DB_USER/" $CONFIG_FILE
fi

if [ ! -z "$DB_PASSWORD" ]; then
    sed -i "s/^password.*/password = $DB_PASSWORD/" $CONFIG_FILE
fi

if [ ! -z "$DB_NAME" ]; then
    sed -i "s/^dbname.*/dbname = $DB_NAME/" $CONFIG_FILE
fi

if [ ! -z "$DB_SSLMODE" ]; then
    sed -i "s/^sslmode.*/sslmode = $DB_SSLMODE/" $CONFIG_FILE
fi

if [ ! -z "$HTTP_LISTENER" ]; then
    sed -i "s/^http_listener.*/http_listener = $HTTP_LISTENER/" $CONFIG_FILE
fi

if [ ! -z "$HTTP_PORT" ]; then
    sed -i "s/^http_port.*/http_port = $HTTP_PORT/" $CONFIG_FILE
fi

cd /app

/app/infinityapi -migrate

exec "$@"
