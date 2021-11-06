#!/usr/bin/env python3

from urllib.parse import urlparse
import pprint
import argparse
import json
import subprocess
import sys
import os


parser = argparse.ArgumentParser()
parser.add_argument("url")
parser.add_argument("--cwd", dest="cwd", default=".")
parser.add_argument("-X", dest="method", default="GET")
parser.add_argument("-H", action="append", dest="headers", default=[])
parser.add_argument("--body", dest="body", default="")
args = parser.parse_args()

os.chdir(args.cwd)

scheme = urlparse(f"http://{args.url}")

headers = {
  s[0]: s[1] for s in map(lambda h: h.split(': ', 1), args.headers)
}

data = {
  "version": "2.0",
  "routeKey": "$default",
  "rawPath": scheme.path,
  "rawQueryString": scheme.query,
  "headers": headers,
  "body": args.body,
  "isBase64Encoded": False,
  "requestContext": {
    "http": {
      "method": args.method,
      "path": scheme.path,
      "protocol": "HTTP/1.1",
      "sourceIP": "127.0.0.1",
      "userAgent": "bebra/0.1"
    },
    "routeKey": "$default",
    "stage": "$default",
  }
}
cmd = ["sls", "invoke", "local", "--function", scheme.netloc, "--data", json.dumps(data)]
with subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=sys.stderr) as proc:
  answer = None
  for line in proc.stdout:
    decoded = line.decode()
    if '"statusCode"' in decoded:
      answer = decoded
    else:
      print(decoded, end='', flush=True)

if answer != None:
  response = json.loads(answer.strip())
  if "body" in response:
    response["body"] = json.loads(response["body"])
  pprint.pprint(response)