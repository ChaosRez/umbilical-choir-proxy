# Umbilical Choir (Reverse) Proxy
This repository is part of the [Umbilical Choir project](https://github.com/ChaosRez/umbilical-choir-core).
For other repositories, see the [Umbilical Choir Agent](https://github.com/ChaosRez/umbilical-choir-core) and [Umbilical Choir Release Manager](https://github.com/ChaosRez/umbilical-choir-release-manager) repositories.
----------------

## go runtime
simply pass the `*.go` and `go.mod` files as Golang function source to be deployed by the agent

## build for python runtime
To cross compile for a linux machine, you need to set target when building
```
GOOS=linux GOARCH=amd64 go build -o binary/_gcp-amd64/umbilical-choir-proxy .
GOOS=linux GOARCH=arm64 go build -o binary/_tinyfaas-arm64/umbilical-choir-proxy .
# legacy:
GOOS=linux GOARCH=arm64 go build -o binary/bash-arm-linux/umbilical-choir-proxy .
GOOS=linux GOARCH=arm64 go build -o binary/python-arm-linux/umbilical-choir-proxy .
GOOS=darwin GOARCH=arm64 go build -o binary/bash-m2/umbilical-choir-proxy .

```

## HTTP header for A/B test
For A/B testing with stateless FaaS functions, use the `X-Function-Version` header to specify the function version. This ensures the user stays on the same version, avoiding probabilistic selection.
Alternatively, a client ID can be used to ensure the same user stays on the same version.

### Why we developed a lightweight metric collector inside the Agent for FaaS functions, and not using the de facto metric collector Prometheus?
Prometheus supports a pull approach, where it pulls the metrics from the target. This is not suitable for FaaS functions as they are short-lived and the metrics will be lost.
There is a [pushgateway](https://github.com/prometheus/pushgateway) which can be used to push the metrics to Prometheus, but it is not suitable for FaaS functions as it will reset the metrics every time the function is called.
- In tinyFaaS, a static counter had to be made in Python as the function stays running in contrast to the binary option which just runs a shell command every time. This is not practical for stateless FaaS functions.
- Prometheus counter doesn't help as it only pulls the absolute value and even in pushgateway it depends on the local counter in your app which will be reset every time in a stateless FaaS environment.
- Prometheus will take 1s samples of pushgateway and NOT all response times. For more, read [When to use the Pushgateway](https://prometheus.io/docs/practices/pushing/)

## Acknowledgement
This repository is part of the [Umbilical Choir](https://github.com/ChaosRez/umbilical-choir-core) project.
If you use this code, please cite our paper and reference this repository in your own code repository.

## Research

If you use any of Umbilical Choir's software components ([Release Manager](https://github.com/ChaosRez/umbilical-choir-release-manager), [Proxy](https://github.com/ChaosRez/umbilical-choir-proxy), and [Agent](https://github.com/ChaosRez/umbilical-choir-core)) in a publication, please cite it as:

### Text

M. Malekabbasi, T. Pfandzelter, and D. Bermbach, **Umbilical Choir: Automated Live Testing for Edge-To-Cloud FaaS Applications**, 2025.

### BibTeX

```bibtex
@inproceedings{malekabbasi2025umbilical,
    author = "Malekabbasi, Mohammadreza and Pfandzelter, Tobias and Bermbach, David",
    title = "Umbilical Choir: Automated Live Testing for Edge-to-Cloud FaaS Applications",
    booktitle = "Proceedings of the 9th IEEE International Conference on Fog and Edge Computing",
    pages = "11--18",
    month = may,
    year = 2025,
    acmid = MISSING,
    publisher = "IEEE",
    address = "New York, NY, USA",
    series = "ICFEC '25",
    location = "Tromso, Norway",
    numpages = MISSING,
    url = "https://doi.org/10.1109/ICFEC65699.2025.00010",
    doi = "10.1109/ICFEC65699.2025.00010"
}
```
