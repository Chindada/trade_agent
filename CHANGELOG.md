# CHANGELOG

## Unreleased (2022-01-19)

### New feature

- **targets**: change target of trade day if no targets, update protoc, no close to log([`8c244ad`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/8c244ad3dc3397f36cfcffb234a3e6d7790c0353)) (@TimHsu@M1BP-20211221)
- **quota**: add check quota in order flow([`61a2c74`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/61a2c74e082f64e8ccb792ed25f0b5e12e8b544f)) (@TimHsu@M1BP-20211221)
- **trade**: add max loss, quota struct([`fa7e186`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/fa7e1861b4f35b3da90f3b2e04dc5880d4a5d5d2)) (@TimHsu@M1BP-20211221)

### Bugs fixed

- **mqtopic**: fix wrong topic string, add max loss, modify config([`2aaf5fa`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/2aaf5fadf48047e6efddad154354dc8e4bf9f77c)) (@TimHsu@M1BP-20211221)

## v1.0.1-alpha (2022-01-13)

### New feature

- **order**: simulation order on sinopac mq srv, fix status no result will panic([`9d83fcf`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/9d83fcf2eb2a449b390af45ab4f89993b6e8363e)) (@TimHsu@M1BP-20211221)
- **cache**: refactor getter, add trade day open end time to avoid do clear order over time([`d5884c9`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/d5884c9645df083ecb6b10604c3c0a56d2c46a21)) (@TimHsu@M1BP-20211221)
- **cache**: split high frequency usage to different type, add volume pr high([`9b611aa`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/9b611aa0a6474fbd740aa1db50e437d8692db558)) (@TimHsu@M1BP-20211221)
- **test**: add runtime path to add test file([`93bee1d`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/93bee1d674c38b6d4aab346c0f66df8dccd2507e)) (@TimHsu@M1BP-20211221)

### Bugs fixed

- **order**: redesign order flow, statsu not finished, remove save analyze, add api status by order id([`00be677`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/00be6771391eb3f54be17c95a75f60c511af6c44)) (@TimHsu@M1BP-20211221)

## v1.0.0-alpha (2022-01-11)

### New feature

- **config**: add check config file value, modify db agent dsn string([`4fdbe15`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/4fdbe15eea36e1e81ca60104f208c495abd86d49)) (@TimHsu@M1BP-20211221)
- **order**: add clear unfinished, split high frequency cache, modify realtime tick analyze([`25f14e7`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/25f14e71ceb5d408a2bba4d9f63b1768f1293db5)) (@TimHsu@M1BP-20211221)
- **anaylze**: config add analyze period min,max, fix outinratio always first tick([`63b30c1`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/63b30c1dc3a790ddacdfcc83655205501d62a46f)) (@TimHsu@M1BP-20211221)
- **cache**: add all cache a setter([`f351749`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/f351749bbb7244e7a6190b0577648d32b98ee119)) (@TimHsu@M1BP-20211221)
- **analyze**: temp save pr outinratio to db([`06faf73`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/06faf7305476a32deeeacb861bf27b6bacdcac0a)) (@TimHsu@M1BP-20211221)
- **analyze**: realtime tick action generator alpha, make sure stock is day trade([`211ce79`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/211ce79bba3fa98690a865551a8ad7d16d882bea)) (@TimHsu@M1BP-20211221)
- **ci**: change from manual to auto deployment([`6e57169`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/6e57169c2815f00b0842398de3fdb9717ddda615)) (@TimHsu@M1BP-20211221)
- **enhancement**: add check port or lock, balance api handlers([`e53e30d`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/e53e30d1483f3ebc55cba480f380567ed71d9025)) (@TimHsu@M1BP-20211221)
- **analyze**: add analyze config, fix total time to nano, realtime target change to every 60 second([`2ad058c`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/2ad058c2c745bab805f4b61b6ebf1257d9147f1b)) (@TimHsu@M1BP-20211221)
- **target**: add every 30 second add target in open time([`788274a`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/788274a9152d3115a1db77b54dd4171f22ab2fc0)) (@TimHsu@M1BP-20211221)
- **tradeday**: add check is open and save to cache([`bef64b0`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/bef64b00dea3c7190991ea7d419724b8b5ca0918)) (@TimHsu@M1BP-20211221)
- **layout**: add health check, move global variable to a struct with getter, setter([`e3b8408`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/e3b8408b4350e2612c8edd436a291f8c231ed190)) (@TimHsu@M1BP-20211221)
- **basepath**: move basepath to global, initial in main package([`55cbf19`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/55cbf19f19480ba7d21ad9c2138bea77ede9ef54)) (@TimHsu@M1BP-20211221)
- **layout**: refactor cache, one type one cache([`26d3570`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/26d3570537c7013f00981055b255306bc3c47ca5)) (@TimHsu@M1BP-20211221)
- **layout**: add trade switch, modify config, biasrate for quantity, order count detail([`832745e`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/832745e6f7558c7074f686f62c5ee9c6a6381643)) (@TimHsu@M1BP-20211221)
- **layout**: finish all basic flow([`c33426c`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/c33426c36e92a4f6cf03c45c07b687c70f7dd25b)) (@TimHsu@M1BP-20211221)
- **layout**: add kbar analyze flow, add history tick analyze method([`e32efd9`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/e32efd9a5d7e0f540f396ed28ea330863e172401)) (@TimHsu@M1BP-20211221)
- **layout**: add order forward, reverse cache, analyze flow([`54c2d8a`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/54c2d8a7b808182f3767b091ce8459a28bf95f5c)) (@TimHsu@M1BP-20211221)
- **layout**: add order flow, tick analyze([`91fa39b`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/91fa39b7fd427bfb0dfc376605d1bd10f09d97b8)) (@TimHsu@M1BP-20211221)
- **layout**: place order v0.1, add subscribe target topic, order add buy later action([`7a38f5d`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/7a38f5d53defa38d7178b511d5c8af864fb6ecbd)) (@TimHsu@M1BP-20211221)
- **layout**: remove simulation, analyzed, add history close, realtime tick, bidask channel cache([`448540e`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/448540e60539f6960e92945e33a6f533aed940c4)) (@TimHsu@M1BP-20211221)
- **layout**: add check all db record, add kbar([`b8d370f`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/b8d370fb3741a854ba2652c70bc5b6e3c719c9bc)) (@TimHsu@M1BP-20211221)
- **layout**: add order frame, remove future, modify history([`fb88059`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/fb8805929e50b59829810e66c9e0c11dd609a1d2)) (@TimHsu@M1BP-20210907)
- **taerror**: new pkg taerror, add unmarshalproto method, sort snapshot([`30de565`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/30de565fdef34d052eb2ed047c12440fe3361eae)) (@TimHsu@M1BP-20210907)
- **layout**: basic func alpha, stock, subscribe, history([`92b439c`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/92b439c93e4f7ffef8c1154b33ae684a3e75d6d6)) (@TimHsu@M1BP-20210907)
- **db**: add all common crud for models([`3e5870a`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/3e5870a36e22224061df7bff27693336a76cc4ee)) (@TimHsu@M1BP-20210907)
- **layout**: add cache lock, cache method for key, initial all model, all db relation([`633caeb`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/633caeb7c7c55ac6633aad162162c1bdfffe0e6b)) (@TimHsu@M1BP-20210907)
- **layout**: change config path method, add trade day, target stock module some feature([`b30ad5a`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/b30ad5a52581cb66fb2efbcc9c5d48bfcf5145e8)) (@TimHsu@M1BP-20210907)
- **layout**: ensure all basic pkg works, sinopacapi remove protobuf([`a9cee17`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/a9cee178ccf7a13fce7f4d612165735d92f1cd6e)) (@TimHsu@M1BP-20210907)
- **log**: add log date folder([`2cf5996`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/2cf5996fe038388fe065f9db9440fc607c00eb71)) (@TimHsu@M1BP-20210907)
- **layout**: add mqhandler, certs, eventbus, yaml config([`8194ccf`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/8194ccf73accce5c5e148b6a9917c98628e53f59)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **config**: fix wrong type in some field, modify config checker, rename some confit struct([`17d9f3d`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/17d9f3d828ce9ed7b202e0bac9d733509ed570aa)) (@TimHsu@M1BP-20211221)
- **order**: fix sleep after continue cause cpu high, save last period volume([`e0909f4`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/e0909f414f46dd6b849840e0792c65678bf8769f)) (@TimHsu@M1BP-20211221)
- **config**: fix missing tick_analyze_period, tick_analyze_max_period([`1cf4bfc`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/1cf4bfcab16603c8b2ea7cd2222ede31497d1105)) (@TimHsu@M1BP-20211221)
- **cache**: fix wrong key in history tick analyze, move simtrade position([`740aa6b`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/740aa6b5f521f914c84dcd08bcd341859e6a2732)) (@TimHsu@M1BP-20211221)
- **ci**: fix missing files copy in dockerfile([`33dd6dc`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/33dd6dcf68adcb827d86fc98d853041ba95e7750)) (@TimHsu@M1BP-20211221)
- **swagger**: fix error remove docs([`bfa8433`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/bfa8433faab5243b28e32112db5220e842d572a2)) (@TimHsu@M1BP-20211221)
- **layout**: rename some easy name, remove wait group in target([`832620b`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/832620bdb12fdc7a48bd03d62e7be71bee17d550)) (@TimHsu@M1BP-20210907)
- **panic**: fix missing file by wrong git clean fxd([`145af44`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/145af44f08dd7df82b6e7a5b752f2e19ac6bce4f)) (@TimHsu@M1BP-20210907)
- **mqhandler**: fix deadlock when subscribe([`6bcd703`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/6bcd7035aaa7e04c06370ad5af9d4267cdb17d1a)) (@TimHsu@M1BP-20210907)
- **mqhandler**: fix format warning on panicf([`b46eee5`](https://gitlab.tocraw.com/trader_v2/trade_agent/commit/b46eee595bab0c385008728f46ac469b0037d584)) (@TimHsu@M1BP-20210907)
