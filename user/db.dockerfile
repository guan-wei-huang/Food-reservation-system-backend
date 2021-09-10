FROM postgres:13.4-alpine
ADD . /docker-entrypoint-initdb.d
