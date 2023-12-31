---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "apisix_route Resource - terraform-provider-apisix"
subcategory: ""
description: |-
  Manages APISIX routes.
---

# apisix_route (Resource)

Manages APISIX routes.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `desc` (String) Description of usage scenarios.
- `enable_websocket` (Boolean) Enables a websocket. Set to `false` by default.
- `filter_func` (String) Matches based on a user-defined filtering function.Used in scenarios requiring complex matching. These functions can accept an input parameter `vars` which can be used to access the Nginx variables.
- `host` (String) Matches with domain names such as `foo.com` or PAN domain names like `*.foo.com`.
- `hosts` (List of String) Matches with any one of the multiple `host`s specified in the form of a non-empty list.
- `labels` (Map of String) Attributes of the Service specified as key-value pairs.
- `methods` (List of String) Matches with the specified HTTP methods. Matches all methods if empty or unspecified.
- `name` (String) Identifier for the route.
- `plugin_config_id` (String) Plugin config bound to the Route.
- `plugins` (String) Plugins that are executed during the request/response cycle.
- `priority` (Number) If different Routes matches to the same `uri`, then the Route is matched based on its `priority`.A higher value corresponds to higher priority.It is set to `0` by default.
- `remote_addr` (String) Matches with the specified IP address in standard IPv4 format (`192.168.1.101`), CIDR format (`192.168.1.0/24`), or in IPv6 format.
- `remote_addrs` (List of String) Matches with any one of the multiple `remote_addrs` specified in the form of a non-empty list.
- `script` (String) Used for writing arbitrary Lua code or directly calling existing plugins to be executed.
- `service_id` (String) Configuration of the bound Service.
- `status` (Number) Enables the current Route. Set to `1` (enabled) by default. `1` to enable, `0` to disable
- `timeout` (Attributes) Sets the timeout (in seconds) for connecting to, and sending and receiving messages to and from the Upstream. (see [below for nested schema](#nestedatt--timeout))
- `upstream_id` (String) Id of the Upstream service.
- `uri` (String) Matches the uri.
- `uris` (List of String) Matches with any one of the multiple `uri`s specified in the form of a non-empty list.
- `vars` (String) Matches based on the specified variables consistent with variables in Nginx. Takes the form `[[var, operator, val], [var, operator, val], ...]]`.

### Read-Only

- `id` (String) Identifier of the route.

<a id="nestedatt--timeout"></a>
### Nested Schema for `timeout`

Required:

- `connect` (Number)
- `read` (Number)
- `send` (Number)

## Import

Import is supported using the following syntax:

```shell
# Route can be imported by specifying the numeric identifier.
terraform import apisix_route.example 123
```
