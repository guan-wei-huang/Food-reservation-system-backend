FROM gangstead/postgis:13-arm
ADD init.sql /docker-entrypoint-initdb.d