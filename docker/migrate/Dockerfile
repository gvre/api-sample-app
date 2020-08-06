FROM alpine:latest

RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.10.0/migrate.linux-amd64.tar.gz && \
    tar xvf migrate.linux-amd64.tar.gz && \
    mv migrate.linux-amd64 /usr/local/bin/migrate

ENV PGHOST $PGHOST
ENV PGUSER $PGUSER
ENV PGPASSWORD $PGPASSWORD
ENV PGDATABASE $PGDATABASE
