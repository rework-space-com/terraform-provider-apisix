# Terraform APISIX Provider
Terraform provider to configure Apache APISIX® using its API.

## APISIX Compatibility
Tested with Apache APISIX® `3.8.0`.

## Usage
The provider configuration method loads configuration data either from environment variables, or from the provider block in Terraform configuration. 

```terraform
provider "apisix" {
  endpoint = "http://127.0.0.1:9180"
  api_key  = "edd1c9f034335f136f87ad84b625c8f1"
}
```
You can use the `APISIX_ENDPOINT` and `APISIX_APIKEY` environment variables for the provider configuration.
```bash
$ APISIX_ENDPOINT=http://127.0.0.1:9180 \
APISIX_API_KEY=edd1c9f034335f136f87ad84b625c8f1 \
terraform plan
```

## Development

-	[Terraform](https://www.terraform.io/downloads.html) 1.x
-	[Go](https://golang.org/doc/install) 1.20.x (to build the provider plugin)

## Contributing to the provider

To contribute, please read the [contribution guidelines](_about/CONTRIBUTING.md). You may also [report an issue](https://github.com/rework-space-com/terraform-provider-apisix/issues/new/choose).
