# Terraform APISIX Provider

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-apisix
```

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

## Install provider locally
You can install provider locally for development and testing.
Use the `go install` command to compile the provider into a binary and install it in your `GOBIN` path.
```shell
$ go install .
```

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

Create a new file called `.terraformrc` in your home directory (`~`), then add the `dev_overrides` block below. Change the `<PATH>` to the value returned from the `go env GOBIN` command.

If the `GOBIN` go environment variable is not set, use the default path, `/home/<Username>/go/bin`.

```terraform
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/apisix" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}

```

### Locally installed provider configuration
Use the folloving provider configration block for the locally installed provider
```terraform
terraform {
  required_providers {
    apisix = {
      source = "hashicorp.com/edu/apisix"
    }
  }
}

provider "apisix" {
  endpoint = "http://127.0.0.1:9180"
  api_key  = "edd1c9f034335f136f87ad84b625c8f1"
}
```

### Start Apache APISIX locally
You can start local APISIX instance using the provided Docker Compose file. The file was adapted from the [apisix-docker repository](https://github.com/apache/apisix-docker/blob/master/example/docker-compose.yml).

In another terminal window, navigate to the `docker_compose` directory.
```bash
cd docker_compose
```
Run `docker-compose up` to spin up a local instance of APISIX on port 9180.
```bash
docker-compose up
```
Leave this process running in your terminal window. In the original terminal window, verify that APISIX is running by sending a request.
```bash
$ curl "http://127.0.0.1:9180/apisix/admin/services/" \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1'
```
The response indicates that apisix is running successfully:
```json
{"total":0,"list":[]}
```
The credentials for the test user are defined in the `docker_compose/apisix_conf/config.yaml` file.