resource "apisix_secret" "example" {
  id = "123"
  gcp = {
    auth_file  = "/path/to/file.json"
    ssl_verify = true
  }
}
