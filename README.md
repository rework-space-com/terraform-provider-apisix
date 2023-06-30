# Terraform APISIX Provider
The APISIX Provider allows Terraform to manage APISIX resources.

## APISIX Compatibility
Tested on APISIX version `3.3.0`

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 1.x
-	[Go](https://golang.org/doc/install) 1.19.x (to build the provider plugin)

## Contributing to the provider

To contribute, please read the [contribution guidelines](_about/CONTRIBUTING.md). You may also [report an issue](https://github.com/rework-space-com/terraform-provider-apisix/issues/new/choose).


## Provider configuration
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
