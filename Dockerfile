FROM docker.io/golang:alpine3.20 as build

ENV CGO_ENABLED=0
ADD . /
WORKDIR /
RUN go build .

FROM docker.io/alpine:3.20
COPY --from=build --chmod=0755 --chown=root:root /supreme-logger-3000 /supreme-logger-3000
USER nobody
ENTRYPOINT [ "/supreme-logger-3000" ]

