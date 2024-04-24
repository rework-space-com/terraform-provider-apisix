## Developing the provider

Thank you for your interest in contributing.

## Documentation

Terraform [provider development documentation](https://www.terraform.io/docs/extend/) provides a good start into developing an understanding of provider development.


## Building the provider

There is a [makefile](../GNUmakefile) to help build the provider. You can build the provider by running `make build`.

```shell
$ make build
```

### Tests
In order to run the full suite of Acceptance tests, run `make testacc`.


```shell
$ make testacc
```

## Install provider locally
You can install provider locally for development and testing.
Use the `make install` command to compile the provider into a binary and install it in your `GOBIN` path.
```shell
$ make install
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

### Create documentation

When creating or updating resources/data resources please make sure to update the examples in the respective folder (`./examples/resources/<name>` for resources, `./examples/data-sources/<name>` for data sources)

Next you can use the following command to generate the terraform documentation from go files

```shell
make doc
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