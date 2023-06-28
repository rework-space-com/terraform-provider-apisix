resource "apisix_consumer_group" "example" {
  id   = "123"
  desc = "Example of the consumer group resource usage"
  plugins = jsonencode(
    {
      prometheus = {
        prefer_name = false
      }
    }
  )
  labels = {
    version = "v1"
  }
}