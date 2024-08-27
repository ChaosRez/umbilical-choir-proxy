#!/usr/bin/env python3

import subprocess
import os

def fn(request):
    """Call the 'umbilical-choir-proxy' binary with the 'input' as an input argument."""
    input = request.data.decode("utf-8")
    print(f"Call {Counter.increment_count()} with input: {input}")
    if input is None:
        input = ""  # Replace None with an empty string if necessary

    result = run_proxy(input)
    if result.returncode != 0:
        print(f"Error running binary proxy: {result.stderr}")
        return f"Error running binary proxy: {result.stderr}"

    return result.stdout

def run_proxy(input: str):
    return subprocess.run(["./umbilical-choir-proxy", input], capture_output=True, text=True)


# NOTE kept this counter only for running chmod once! otherwise it is not needed anymore
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