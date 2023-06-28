resource "apisix_consumer" "example" {
  username = "example"
  desc     = "Example of the consumer"
  labels = {
    "version" = "v1"
  }
  plugins = jsonencode(
    {
      basic-auth = {
        username = "example"
        password = "changeme2"
      }
    }
  )
}
