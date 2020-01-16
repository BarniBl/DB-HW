FROM golang:1.13-stretch AS builder

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -v -work

FROM ubuntu:19.04
ENV DEBIAN_FRONTEND=noninteractive
ENV PGVER 11
ENV PORT 5000
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432
ENV POSTGRES_DB forum
ENV POSTGRES_USER forum
ENV POSTGRES_PASSWORD 7396
EXPOSE $PORT

RUN apt-get update && apt-get install -y postgresql-$PGVER

USER postgres

RUN service postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD '7396';" &&\
    createdb -O forum forum &&\
    service postgresql stop

#COPY config/pg_hba.conf /etc/postgresql/$PGVER/main/pg_hba.conf
#COPY config/postgresql.conf /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

COPY dum_hw_pdb.sql .
RUN ls /usr
COPY --from=builder /usr/src/app/DB-HW .
CMD service postgresql start && ./DB-HW
