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

## HTTP header for A/B test
For A/B testing with stateless FaaS functions, use the `X-Function-Version` header to specify the function version. This ensures the user stays on the same version, avoiding probabilistic selection.

## Why not prometheus for FaaS functions
Prometheus supports a pull approach, where it pulls the metrics from the target. This is not suitable for FaaS functions as they are short lived and the metrics will be lost.
There is a pushgateway which can be used to push the metrics to prometheus, but it is not suitable for FaaS functions as it will reset the metrics every time the function is called.
- had to make static counter in python as the function stays running in contrast to binary option which just runs a shell command everytime
- promethus counter doesn't help as it only pulls the value and even in pushgateway it depends on the local value in your app which will reset everytime
- prometheus will take 1s samples of pushgateway and NOT all response times