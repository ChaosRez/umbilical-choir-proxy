# Umbilical Choir (Reverse) Proxy
## build
To cross compile for a linux machine, you need to set target either as environment param:  
```
export GOOS=linux GOARCH=arm64
```
or directly when building
```
GOOS=linux GOARCH=amd64 go build -o binary/_gcp-amd64/umbilical-choir-proxy .
GOOS=linux GOARCH=arm64 go build -o binary/_tinyfaas-arm64/umbilical-choir-proxy .

GOOS=linux GOARCH=arm64 go build -o binary/bash-arm-linux/umbilical-choir-proxy .
GOOS=linux GOARCH=arm64 go build -o binary/python-arm-linux/umbilical-choir-proxy .
GOOS=darwin GOARCH=arm64 go build -o binary/bash-m2/umbilical-choir-proxy .

```