resource "apisix_ssl_certificate" "example" {
  certificate = file("example.crt")
  private_key = file("example.key")
  type        = "server"
  labels = {
    "version" = "v1"
  }
}