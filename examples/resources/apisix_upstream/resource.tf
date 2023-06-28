resource "apisix_upstream" "example" {
  name = "Example"
  desc = "Example of the upstream resource usage"
  type = "roundrobin"
  labels = {
    version = "v1"
  }
  nodes = [
    {
      host   = "127.0.0.1"
      port   = 1980
      weight = 1
    },
    {
      host   = "127.0.0.1"
      port   = 1970
      weight = 1
    },
  ]
  keepalive_pool = {
    idle_timeout = 5
    requests     = 10
    size         = 15
  }
  checks = {
    active = {
      host      = "example.com"
      port      = 8888
      timeout   = 5
      http_path = "/status"
      healthy = {
        interval  = 2,
        successes = 1
      }
      unhealthy = {
        interval      = 1
        http_failures = 2
      }
    }
    passive = {
      healthy = {
        http_statuses = [200, 201]
      }
      unhealthy = {
        http_statuses = [500]
        http_failures = 3
        tcp_failures  = 3
      }
    }
  }
}