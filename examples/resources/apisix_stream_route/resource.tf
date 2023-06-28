resource "apisix_upstream" "example" {
  type = "roundrobin"
  nodes = [
    {
      host   = "127.0.0.1"
      port   = 8080
      weight = 1
    }
  ]
}

resource "apisix_stream_route" "example" {
  upstream_id = apisix_upstream.example.id
  remote_addr = "127.0.0.1"
  server_addr = "127.0.0.1"
  server_port = 8080
  sni         = "example.com"
}