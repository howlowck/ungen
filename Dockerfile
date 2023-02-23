FROM alpine:3.17

RUN apk add zip

COPY ./ungen /usr/local/bin/

ENTRYPOINT ["ungen"]
