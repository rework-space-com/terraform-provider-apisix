resource "apisix_plugin_config" "example" {
  id   = "123"
  desc = "Example of the plugin config resource usage"
  plugins = jsonencode(
    {
      prometheus = {
        prefer_name = true
      }
    }
  )
  labels = {
    version = "v1"
  }
}