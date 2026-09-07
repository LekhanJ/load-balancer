from http.server import HTTPServer, BaseHTTPRequestHandler
import os

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        port = os.environ.get("PORT", "unknown")
        
        self.send_response(200)
        self.end_headers()
        
        self.wfile.write(f"Response from Python server {port}\n".encode())
    
    def do_HEAD(self):
        self.send_response(200)
        self.end_headers()

port = int(os.environ.get("PORT", "5001"))
server = HTTPServer(("localhost", port), Handler)
print(f"Starting Python backend on port {port}")

server.serve_forever()