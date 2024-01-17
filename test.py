#!/usr/bin/env python3

import urllib.parse

import httpx

TEST_KEY = """this
has
new
lines
"""
VALUE = "foo"


def main(client: httpx.Client) -> None:
    encoded_key = urllib.parse.quote(TEST_KEY)
    # put
    encoded_url = f"/kv/v1/put/{encoded_key}"
    print(f"encoded put url: {encoded_url}")
    resp = client.put(
        url=encoded_url,
        json={"value": VALUE}
    )
    resp.raise_for_status()

    # get
    encoded_url = f"/kv/v1/get/{encoded_key}"
    print(f"encoded get url: {encoded_url}")
    resp = client.get(
        url=encoded_url,
    )
    resp.raise_for_status()
    data = resp.json()
    assert data["value"] == VALUE


if __name__ == '__main__':
    with httpx.Client(base_url="http://localhost:8081") as http:
        main(client=http)
