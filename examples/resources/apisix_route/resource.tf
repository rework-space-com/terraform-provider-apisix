resource "apisix_route" "example" {
  name         = "Example"
  desc         = "Example of the route configuration"
  uris         = ["/api/v1", "/status"]
  hosts        = ["foo.com", "*.bar.com"]
  remote_addrs = ["10.0.0.0/8"]
  methods      = ["GET", "POST"]
  vars = jsonencode(
    [["http_user", "==", "ios"]]
  )
  timeout = {
    connect = 3
    send    = 3
    read    = 3
  }
  plugins = jsonencode(
    {
      ip-restriction = {
        blacklist = ["10.10.10.0/24"]
        message   = "Access denied"
      }
    }
  )
  labels = {
    "version" = "v1"
  }
}
