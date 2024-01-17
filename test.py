#!/usr/bin/env python3

import urllib.parse

import httpx

TEST_KEY = """this
has
new
lines
"""
VALUE = "foo"


def main(*, client: httpx.Client, url: str) -> None:
    resp = client.put(url=url, json={"value": VALUE})
    resp.raise_for_status()

    # get
    resp = client.get(url=url)
    resp.raise_for_status()
    data = resp.json()
    assert data["value"] == VALUE


if __name__ == "__main__":
    with httpx.Client(base_url="http://localhost:8081") as http:
        main(client=http, url=f"/kv/v1/{urllib.parse.quote(TEST_KEY)}")
