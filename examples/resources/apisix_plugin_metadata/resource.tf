resource "apisix_plugin_metadata" "syslog" {
  id = "syslog"
  metadata = jsonencode(
    {
      "log_format" = {
        "@timestamp" = "$time_iso8601"
        "client_ip"  = "$remote_addr"
        "host"       = "$host"
      }
    }
  )
}