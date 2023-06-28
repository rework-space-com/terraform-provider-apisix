resource "apisix_global_rule" "example" {
  id = "123"
  plugins = jsonencode(
    {
      prometheus = {
        prefer_name = true
      }
    }
  )
}