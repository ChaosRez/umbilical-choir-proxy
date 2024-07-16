#!/usr/bin/env python3

import typing
import subprocess
import os
import random
import time
from prometheus_client import CollectorRegistry, Gauge, push_to_gateway
import requests

# Initialize a new CollectorRegistry
registry = CollectorRegistry()
# Initialize Prometheus metrics
count_p = Gauge('call_count', 'Number of calls to fn', registry=registry)
response_time_1 = Gauge('f1_time', 'Response time for the f1 (ms)', registry=registry)
response_time_2 = Gauge('f2_time', 'Response time for the f2 (ms)', registry=registry)
proxy_time = Gauge('proxy_time', 'Total proxy runtime (ms)', registry=registry)

# Environment variables
host = os.getenv("HOST")
port = os.getenv("PORT")
f1_name = os.getenv("F1NAME")
f2_name = os.getenv("F2NAME")
program = os.getenv("PROGRAM")

# Global http client session
session = requests.Session()

def fn(input: typing.Optional[str]) -> typing.Optional[str]:  # TODO: add async I/O for function and metrics. Won't need all samples
    start_proxy = time.time()
    print(f"Call {Counter.increment_count()} with input: {input}")
    if input is None:
        input = ""  # Replace None with an empty string if necessary

    # Call one of the function versions
    resp, elap, is_f2 = uniform_random_call(input)
    if resp is None:
        print("Error running uniformCallAndLog")
        return

    # Print stdout
    if is_f2:
        print(f"resp (f2): {resp} \n took: {elap}")
    else:
        print(f"resp (f1): {resp} \n took: {elap}")

    # Update metrics
    elap_total = time.time() - start_proxy
    proxy_time.set(elap_total * 1000)  # Convert to milliseconds
    count_p.set(Counter.get_count())
    push_metrics()
    return resp


def uniform_random_call(input_arg):
    choice = random.randint(0, 1)  # Randomly choose between f1 and f2

    if choice == 0:
        return f1_call(input_arg)
    else:
        return f2_call(input_arg)

def f1_call(input_arg):
    try:
        start = time.time()
        resp = session.post(f"http://{host}:{port}/{f1_name}", data=input_arg)
        elap = time.time() - start
        response_time_1.set(elap * 1000)  # Convert to milliseconds
        return resp.text, elap, False
    except requests.ConnectionError as e:
        print(f"Connection error: {e}")
        return "Error", 0, False

def f2_call(input_arg):
    try:
        start = time.time()
        resp = session.post(f"http://{host}:{port}/{f2_name}", data=input_arg)
        elap = time.time() - start
        response_time_2.set(elap * 1000)  # Convert to milliseconds
        return resp.text, elap, True
    except requests.ConnectionError as e:
        print(f"Connection error: {e}")
        return "Error", 0, False

def push_metrics(*metrics):
    # Push all metrics in the registry to the Pushgateway
    push_to_gateway(f'{host}:9091', job='umbilical-choir', registry=registry,
                    grouping_key={'program': program})

class Counter:
    count = None
    first_call = True

    @staticmethod
    def get_count():
        if Counter.count is None:  # memoize
            Counter.count = 0
        return Counter.count

    @staticmethod
    def increment_count():
        if Counter.count is None:
            Counter.count = 0
        Counter.count += 1
        return Counter.count