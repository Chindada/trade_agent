#!/bin/bash

pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile stop
rm -rf ./data/trade_agent
mkdir -p ./data/trade_agent

initdb ./data/trade_agent
pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile start

echo "\du
CREATE ROLE postgres WITH LOGIN PASSWORD 'asdf0000';
ALTER USER postgres WITH SUPERUSER;
\du" > sql_script

psql postgres -f sql_script
rm -rf sql_script

pg_ctl -D ./data/trade_agent -l ./data/trade_agent/logfile stop
