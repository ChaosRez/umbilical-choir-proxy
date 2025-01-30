import typing
import subprocess
import os
import time

def fn(input: typing.Optional[str], headers: typing.Optional[typing.Dict[str, str]]) -> typing.Optional[str]:
    """Call the 'umbilical-choir-proxy' binary with the 'input' as an input argument."""
    start_time = time.time()
    print(f"Call {Counter.increment_count()} with input: {input}")
    if input is None:
        input = ""  # Replace None with an empty string if necessary

    function_choice = headers.get("X-Function-Choice", "") if headers else ""

    start_run_proxy = time.time()
    result = run_proxy(input, function_choice)
    proxy_time = (time.time() - start_run_proxy) * 1000  # Convert to milliseconds

    if result.returncode != 0:
        error_message = f"Error running binary proxy: {result.stderr}"
        return f"{error_message}\nproxy total time: {proxy_time:.2f} ms"

    return f"{result.stdout}\nproxy total time: {proxy_time:.2f} ms"

def run_proxy(input: str, function_choice: str):
    return subprocess.run(["./umbilical-choir-proxy", input, function_choice], capture_output=True, text=True)


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