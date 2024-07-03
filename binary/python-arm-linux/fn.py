#!/usr/bin/env python3

import typing
import subprocess
import os
from concurrent.futures import ThreadPoolExecutor, as_completed
from prometheus_client import CollectorRegistry, Gauge, push_to_gateway

def fn(input: typing.Optional[str]) -> typing.Optional[str]:
    """Call the 'umbilical-choir-proxy' binary with the 'input' as an input argument."""
    print(f"Call {Counter.increment_count()} with input: {input}")
    if input is None:
        input = ""  # Replace None with an empty string if necessary
    with ThreadPoolExecutor(max_workers=2) as executor:
        # Push the counter to the Pushgateway in a new thread
        future_push = executor.submit(push_to_pushgateway, Counter.get_count())

        future_proxy = executor.submit(run_proxy, input)
        for future in as_completed([future_push, future_proxy]):
            if future == future_proxy:
                result = future.result()
                if result.returncode != 0:
                    print(f"Error running binary proxy: {result.stderr}")
                    return f"Error running binary proxy: {result.stderr}"

                return result.stdout

def run_proxy(input: str):
    return subprocess.run(["./umbilical-choir-proxy", input], capture_output=True, text=True)

def push_to_pushgateway(count):
    registry = CollectorRegistry()
    g = Gauge('call_count', 'Number of calls to fn', registry=registry)
    g.set(count)

    push_to_gateway(f'{os.getenv("HOST")}:9091', job='umbilical-choir', registry=registry,
        grouping_key={'program': os.getenv("PROGRAM")})

class Counter:
    count = None
    first_call = True

    @staticmethod
    def get_count():
        if Counter.count is None:  # memoize
            Counter.count = 0
            subprocess.run(["chmod", "755", "umbilical-choir-proxy"]) # only first call
        return Counter.count

    @staticmethod
    def increment_count():
        if Counter.count is None:
            Counter.count = 0
            subprocess.run(["chmod", "755", "umbilical-choir-proxy"]) # only first call
        Counter.count += 1
        return Counter.count