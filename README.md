# TRADE AGENT

[![pipeline status](https://gitlab.tocraw.com/root/trade_agent/badges/main/pipeline.svg)](https://gitlab.tocraw.com/root/trade_agent/-/commits/main)
[![coverage report](https://gitlab.tocraw.com/root/trade_agent/badges/main/coverage.svg)](https://gitlab.tocraw.com/root/trade_agent/-/commits/main)
[![Maintained](https://img.shields.io/badge/Maintained-yes-green)](https://gitlab.tocraw.com/root/trade_agent)
[![Go](https://img.shields.io/badge/Go-1.17.7-blue?logo=go&logoColor=blue)](https://golang.org)
[![OS](https://img.shields.io/badge/OS-Linux-orange?logo=linux&logoColor=orange)](https://www.linux.org/)
[![Virtualization](https://img.shields.io/badge/Container-Docker-blue?logo=docker&logoColor=blue)](https://www.docker.com/)

## Features

[API Docs](http://trade-agent.tocraw.com:16670/swagger/index.html)

### Module Features@2022.01.03

#### Init

- 確認執行檔路徑，並儲存於 `global-Setting-BasePath`
- 由環境變數 `DEPLOYMENT` 確認是否開發環境
  - `docker` 則為開發環境

#### Tasks

- 由 config 中讀取 cron 初始化以下 tasks
  - `CleanEvent`
  - `RestartSinopac`

#### SinopacAPI

- 與 `sinopac_mq_srv` 通訊
- 由 config 中讀取 host, port，並初始化 `AskSinpacMQSRVConnectMQ`, `FetchServerToken`

#### Routers

- API Server

#### MQHandler

- 由 config 中讀取 `mqtt` host, port, username, password
- 建立 `trade_agent_mqtt` 連線
- 無法建立連線則會 lock

#### TradeDay

#### TickProcess

#### Targets

#### Subscribe

#### Stock

#### Order

#### History

#### HealthCheck

#### CloudEvent

#### DBAgent

### Order Module

![callvis](./assets/callvis.svg "callvis")

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)
