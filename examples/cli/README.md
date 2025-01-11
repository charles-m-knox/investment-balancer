# Investment Balancer (CLI)

Copy `config.example.json` to `config.json` and then run:

```bash
make build && make run
```

On Linux you can then do

```bash
xdg-open outputs/investments.csv
```

## Production-ready build

Will produce a file called `investment-balancer-v3-compressed`.

```bash
# install upx for compressing binaries:
# https://upx.github.io/
sudo dnf install upx

make build-prod && make compress-prod
```
