FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y zip

COPY ./ungen /usr/local/bin/

ENTRYPOINT ["ungen"]
