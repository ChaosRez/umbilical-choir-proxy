# Umbilical Choir (Reverse) Proxy
## build
To cross compile for a linux machine, you need to set target either as environment param:  
```
export GOOS=linux GOARCH=arm64
```
or directly when building
```
GOOS=linux GOARCH=arm64 go build -o uc-proxy-linux-arm .
GOOS=darwin GOARCH=arm64 go build -o uc-proxy-mac-m2 .

```