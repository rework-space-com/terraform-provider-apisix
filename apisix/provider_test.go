package apisix

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the APISIX client is properly configured.
	// It is also possible to use the APISIX_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "apisix" {
	endpoint = "http://127.0.0.1:9180"
	api_key  = "edd1c9f034335f136f87ad84b625c8f1"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"apisix": providerserver.NewProtocol6WithError(New("test")()),
	}
)
