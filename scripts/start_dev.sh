#!/bin/bash

pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile stop
pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile start

pkill mosquitto
mosquitto -c /Users/timhsu/dev_projects/mosquitto/trade_agent_mqtt/configs/local_conf.conf
