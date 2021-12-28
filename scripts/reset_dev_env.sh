#!/bin/bash

pkill mosquitto
pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile stop
rm -rf ./data/trade_agent
mkdir -p ./data/trade_agent

initdb ./data/trade_agent

gsed -i "$ a host    all    all    172.20.10.0/24    trust" ./data/trade_agent/pg_hba.conf
gsed -i "$ a listen_addresses = '*'" ./data/trade_agent/postgresql.conf

pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile start

echo "\du
CREATE ROLE postgres WITH LOGIN PASSWORD 'asdf0000';
ALTER USER postgres WITH SUPERUSER;
\du" > sql_script

psql postgres -f sql_script
rm -rf sql_script

# pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile stop

mosquitto -c /Users/timhsu/dev_projects/mosquitto/trade_agent_mqtt/configs/local_conf.conf
