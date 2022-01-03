# CHANGELOG

## Unreleased (2022-01-03)

### New feature

- **analyze**: add analyze config, fix total time to nano, realtime target change to every 60 second([`2ad058c`](https://gitlab.tocraw.com/root/trade_agent/commit/2ad058c2c745bab805f4b61b6ebf1257d9147f1b)) (@TimHsu@M1BP-20211221)
- **target**: add every 30 second add target in open time([`788274a`](https://gitlab.tocraw.com/root/trade_agent/commit/788274a9152d3115a1db77b54dd4171f22ab2fc0)) (@TimHsu@M1BP-20211221)
- **tradeday**: add check is open and save to cache([`bef64b0`](https://gitlab.tocraw.com/root/trade_agent/commit/bef64b00dea3c7190991ea7d419724b8b5ca0918)) (@TimHsu@M1BP-20211221)
- **layout**: add health check, move global variable to a struct with getter, setter([`e3b8408`](https://gitlab.tocraw.com/root/trade_agent/commit/e3b8408b4350e2612c8edd436a291f8c231ed190)) (@TimHsu@M1BP-20211221)
- **basepath**: move basepath to global, initial in main package([`55cbf19`](https://gitlab.tocraw.com/root/trade_agent/commit/55cbf19f19480ba7d21ad9c2138bea77ede9ef54)) (@TimHsu@M1BP-20211221)
- **layout**: refactor cache, one type one cache([`26d3570`](https://gitlab.tocraw.com/root/trade_agent/commit/26d3570537c7013f00981055b255306bc3c47ca5)) (@TimHsu@M1BP-20211221)
- **layout**: add trade switch, modify config, biasrate for quantity, order count detail([`832745e`](https://gitlab.tocraw.com/root/trade_agent/commit/832745e6f7558c7074f686f62c5ee9c6a6381643)) (@TimHsu@M1BP-20211221)
- **layout**: finish all basic flow([`c33426c`](https://gitlab.tocraw.com/root/trade_agent/commit/c33426c36e92a4f6cf03c45c07b687c70f7dd25b)) (@TimHsu@M1BP-20211221)
- **layout**: add kbar analyze flow, add history tick analyze method([`e32efd9`](https://gitlab.tocraw.com/root/trade_agent/commit/e32efd9a5d7e0f540f396ed28ea330863e172401)) (@TimHsu@M1BP-20211221)
- **layout**: add order forward, reverse cache, analyze flow([`54c2d8a`](https://gitlab.tocraw.com/root/trade_agent/commit/54c2d8a7b808182f3767b091ce8459a28bf95f5c)) (@TimHsu@M1BP-20211221)
- **layout**: add order flow, tick analyze([`91fa39b`](https://gitlab.tocraw.com/root/trade_agent/commit/91fa39b7fd427bfb0dfc376605d1bd10f09d97b8)) (@TimHsu@M1BP-20211221)
- **layout**: place order v0.1, add subscribe target topic, order add buy later action([`7a38f5d`](https://gitlab.tocraw.com/root/trade_agent/commit/7a38f5d53defa38d7178b511d5c8af864fb6ecbd)) (@TimHsu@M1BP-20211221)
- **layout**: remove simulation, analyzed, add history close, realtime tick, bidask channel cache([`448540e`](https://gitlab.tocraw.com/root/trade_agent/commit/448540e60539f6960e92945e33a6f533aed940c4)) (@TimHsu@M1BP-20211221)
- **layout**: add check all db record, add kbar([`b8d370f`](https://gitlab.tocraw.com/root/trade_agent/commit/b8d370fb3741a854ba2652c70bc5b6e3c719c9bc)) (@TimHsu@M1BP-20211221)
- **layout**: add order frame, remove future, modify history([`fb88059`](https://gitlab.tocraw.com/root/trade_agent/commit/fb8805929e50b59829810e66c9e0c11dd609a1d2)) (@TimHsu@M1BP-20210907)
- **taerror**: new pkg taerror, add unmarshalproto method, sort snapshot([`30de565`](https://gitlab.tocraw.com/root/trade_agent/commit/30de565fdef34d052eb2ed047c12440fe3361eae)) (@TimHsu@M1BP-20210907)
- **layout**: basic func alpha, stock, subscribe, history([`92b439c`](https://gitlab.tocraw.com/root/trade_agent/commit/92b439c93e4f7ffef8c1154b33ae684a3e75d6d6)) (@TimHsu@M1BP-20210907)
- **db**: add all common crud for models([`3e5870a`](https://gitlab.tocraw.com/root/trade_agent/commit/3e5870a36e22224061df7bff27693336a76cc4ee)) (@TimHsu@M1BP-20210907)
- **layout**: add cache lock, cache method for key, initial all model, all db relation([`633caeb`](https://gitlab.tocraw.com/root/trade_agent/commit/633caeb7c7c55ac6633aad162162c1bdfffe0e6b)) (@TimHsu@M1BP-20210907)
- **layout**: change config path method, add trade day, target stock module some feature([`b30ad5a`](https://gitlab.tocraw.com/root/trade_agent/commit/b30ad5a52581cb66fb2efbcc9c5d48bfcf5145e8)) (@TimHsu@M1BP-20210907)
- **layout**: ensure all basic pkg works, sinopacapi remove protobuf([`a9cee17`](https://gitlab.tocraw.com/root/trade_agent/commit/a9cee178ccf7a13fce7f4d612165735d92f1cd6e)) (@TimHsu@M1BP-20210907)
- **log**: add log date folder([`2cf5996`](https://gitlab.tocraw.com/root/trade_agent/commit/2cf5996fe038388fe065f9db9440fc607c00eb71)) (@TimHsu@M1BP-20210907)
- **layout**: add mqhandler, certs, eventbus, yaml config([`8194ccf`](https://gitlab.tocraw.com/root/trade_agent/commit/8194ccf73accce5c5e148b6a9917c98628e53f59)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **swagger**: fix error remove docs([`bfa8433`](https://gitlab.tocraw.com/root/trade_agent/commit/bfa8433faab5243b28e32112db5220e842d572a2)) (@TimHsu@M1BP-20211221)
- **layout**: rename some easy name, remove wait group in target([`832620b`](https://gitlab.tocraw.com/root/trade_agent/commit/832620bdb12fdc7a48bd03d62e7be71bee17d550)) (@TimHsu@M1BP-20210907)
- **panic**: fix missing file by wrong git clean fxd([`145af44`](https://gitlab.tocraw.com/root/trade_agent/commit/145af44f08dd7df82b6e7a5b752f2e19ac6bce4f)) (@TimHsu@M1BP-20210907)
- **mqhandler**: fix deadlock when subscribe([`6bcd703`](https://gitlab.tocraw.com/root/trade_agent/commit/6bcd7035aaa7e04c06370ad5af9d4267cdb17d1a)) (@TimHsu@M1BP-20210907)
- **mqhandler**: fix format warning on panicf([`b46eee5`](https://gitlab.tocraw.com/root/trade_agent/commit/b46eee595bab0c385008728f46ac469b0037d584)) (@TimHsu@M1BP-20210907)
