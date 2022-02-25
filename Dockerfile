# build-stage
FROM golang:1.17.7-bullseye as build-stage
USER root

ENV GO111MODULE="on"
ENV TZ=Asia/Taipei

WORKDIR /
RUN mkdir build_space
WORKDIR /build_space
COPY . .
WORKDIR /build_space/cmd
RUN go build -o trade_agent

# production-stage
FROM debian:bullseye as production-stage
USER root

ENV DEPLOYMENT=docker
ENV TZ=Asia/Taipei

WORKDIR /
RUN apt update -y && \
    apt install -y tzdata && \
    apt autoremove -y && \
    apt clean && \
    mkdir trade_agent && \
    mkdir trade_agent/configs && \
    mkdir trade_agent/configs/certs && \
    mkdir trade_agent/logs && \
    mkdir trade_agent/scripts && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /trade_agent

COPY --from=build-stage /build_space/cmd/trade_agent ./trade_agent
COPY --from=build-stage /build_space/configs/config.yaml ./configs/config.yaml
COPY --from=build-stage /build_space/configs/certs ./configs/certs
COPY --from=build-stage /build_space/scripts/docker-entrypoint.sh ./scripts/docker-entrypoint.sh

ENTRYPOINT ["/trade_agent/scripts/docker-entrypoint.sh"]
