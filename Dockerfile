# docker-chown | (c) 2020 Icinga GmbH | GPLv2+

FROM golang:buster as build

ADD . /docker-chown
WORKDIR /docker-chown
RUN ["go", "build", "."]


FROM scratch

COPY --from=build /docker-chown/docker-chown /docker-chown
