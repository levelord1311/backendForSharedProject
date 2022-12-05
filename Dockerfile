FROM golang:1.19-alpine3.17 as builder

WORKDIR /usr/local/go/src/

ADD app/ /usr/local/go/src/

# RUN go clean --modchache
RUN go build -mod=readonly -o app cmd/apiserver/main.go

FROM alpine:3.17

COPY --from=builder /usr/local/go/src/app /
COPY --from=builder /usr/local/go/src/config.yml /

CMD [ "/app" ]
