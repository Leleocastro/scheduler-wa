FROM kong:latest

USER root

RUN apt-get update && \
    apt-get install -y build-essential libssl-dev openssl curl

RUN luarocks install luaossl && \
    luarocks install kong-jwt-firebase

USER kong

COPY /config/kong.conf /etc/kong/kong.conf