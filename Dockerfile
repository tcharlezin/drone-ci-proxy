FROM alpine:latest

RUN mkdir /app

COPY drone-ci-proxy /app

CMD [ "/app/drone-ci-proxy" ]