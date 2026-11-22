#!/usr/bin/env python3
"""Heartbeat server: receives POST /report and logs each request (JSON lines)."""
import json
import logging
import os
import sys
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from logging.handlers import RotatingFileHandler

HOST = os.environ.get("HEARTBEAT_HOST", "0.0.0.0")
PORT = int(os.environ.get("HEARTBEAT_PORT", "8080"))
TOKEN = os.environ.get("HEARTBEAT_TOKEN", "")
LOG_FILE = os.environ.get("HEARTBEAT_LOG_FILE", "/heartbeat-requests.log")
MAX_BODY = 64 * 1024

req_log = logging.getLogger("requests")
req_log.setLevel(logging.INFO)
req_log.propagate = False
req_log.addHandler(RotatingFileHandler(LOG_FILE, maxBytes=10_000_000, backupCount=5))
req_log.addHandler(logging.StreamHandler(sys.stdout))
for h in req_log.handlers:
    h.setFormatter(logging.Formatter("%(asctime)s %(message)s"))


class Handler(BaseHTTPRequestHandler):
    def _reply(self, code, body=b"ok\n"):
        self.send_response(code)
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        self._reply(200 if self.path == "/health" else 404)

    def do_POST(self):
        if self.path != "/report":
            return self._reply(404, b"not found\n")
        if TOKEN and self.headers.get("Authorization") != f"Bearer {TOKEN}":
            req_log.warning(json.dumps({"client": self.client_address[0],
                                        "event": "unauthorized"}))
            return self._reply(401, b"unauthorized\n")
        try:
            n = int(self.headers.get("Content-Length", "0"))
            if not 0 < n <= MAX_BODY:
                raise ValueError("bad length")
            data = json.loads(self.rfile.read(n))
        except ValueError:
            return self._reply(400, b"bad request\n")
        req_log.info(json.dumps({"client": self.client_address[0], "data": data}))
        self._reply(200)

    def log_message(self, *args):  # silence default stderr access log
        pass


if __name__ == "__main__":
    srv = ThreadingHTTPServer((HOST, PORT), Handler)
    print(f"listening on {HOST}:{PORT}, log={LOG_FILE}", flush=True)
    try:
        srv.serve_forever()
    except KeyboardInterrupt:
        pass
