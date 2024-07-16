import os
import random
import time
from prometheus_client import Gauge, push_to_gateway, CollectorRegistry
import requests

# Initialize a new CollectorRegistry
registry = CollectorRegistry()

# Initialize Prometheus metrics
response_time_1 = Gauge('response_time_1', 'Response time for the f1 (ms)', registry=registry)
response_time_2 = Gauge('response_time_2', 'Response time for the f2 (ms)', registry=registry)
proxy_time = Gauge('proxy_time', 'Total proxy runtime (ms)', registry=registry)

# Environment variables
host = os.getenv("HOST")
port = os.getenv("PORT")
f1_name = os.getenv("F1NAME")
f2_name = os.getenv("F2NAME")
program = os.getenv("PROGRAM")

def run(input_arg):
    start_proxy = time.time()

    # Call one of the function versions
    resp, elap, is_f2 = uniform_call_and_log(input_arg)
    if resp is None:
        print("Error running uniformCallAndLog")
        return

    # Print stdout
    print(f"resp: {resp} \n took: {elap}")

    # Push total proxy time
    elap_total = time.time() - start_proxy
    proxy_time.set(elap_total * 1000)  # Convert to milliseconds

    # Push metric(s)
    if is_f2:  # f2 was called and not f1
        push_metrics(response_time_2, proxy_time)
    else:
        push_metrics(response_time_1, proxy_time)

def uniform_call_and_log(input_arg):
    choice = random.randint(0, 1)  # Randomly choose between f1 and f2

    if choice == 0:
        return f1_call(input_arg)
    else:
        return f2_call(input_arg)

def f1_call(input_arg):
    start = time.time()
    resp = requests.post(f"http://{host}:{port}/{f1_name}", data=input_arg)
    elap = time.time() - start
    response_time_1.set(elap * 1000)  # Convert to milliseconds
    return resp.text, elap, False

def f2_call(input_arg):
    start = time.time()
    resp = requests.post(f"http://{host}:{port}/{f2_name}", data=input_arg)
    elap = time.time() - start
    response_time_2.set(elap * 1000)  # Convert to milliseconds
    return resp.text, elap, True

def push_metrics(*metrics):
    # Push all metrics in the registry to the Pushgateway
    push_to_gateway(f'{host}:9091', job=program, registry=registry)