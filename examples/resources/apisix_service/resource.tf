resource "apisix_service" "example" {
  name  = "Example"
  hosts = ["foo.com", "*.bar.com"]
  labels = {
    "version" = "v1"
  }
  enable_websocket = true
  plugins = jsonencode(
    {
      limit-count = {
        count                   = 10
        key                     = "remote_addr"
        rejected_code           = 503
        show_limit_quota_header = true
        time_window             = 12
      },
      prometheus = {}
    }
  )
}