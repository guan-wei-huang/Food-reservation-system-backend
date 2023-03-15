FROM postgres:15-alpine
ADD init.sql /docker-entrypoint-initdb.d
