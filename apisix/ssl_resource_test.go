package apisix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSSLResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_ssl_certificate" "test" {
	certificate = "-----BEGIN CERTIFICATE-----\nMIIFPjCCAyagAwIBAgIUe/c2A8XOSdHw60+OrCYfh3dzmm4wDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wHhcNMjMwNjE0MTAxOTMzWhcNMzMw\nNjExMTAxOTMzWjAWMRQwEgYDVQQDDAtleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcN\nAQEBBQADggIPADCCAgoCggIBALAxaEUndKUJ8whS48/pLT5Md/2LuQaov/0lpRcS\naPuteQq0lrbymZA9K0aAUt4Slx1zkKmfW4MKMHtsT5aq3uE7h2+I7v5mnT0EW6+1\nMXwmcLuqQrCVg55ogYHNcoLtsiSpd+vQU+F+MoV0iJpcmf0HcFKN1U9UNldkggl2\nQxIsvaPWxzEZMerHmgNVkAG7YlEXO20bgJfsbQeoBb0n8bNbr09nMcqlcV4hD0H+\nRKf3Xo4o9JdeIJGXbFf8JGolqkFSNUhlwPTkf4Ksndq0nALE0rqAVH848mTdZf58\n+JupCBCvwrajRGI4XKBYziAUe+0GVWIDSalzw6z7OSodOtLQvAIgMauGfJugLHRy\nF677CZIykVv4veAZat7RhI9XeHK83PMxYk+dQvAMOVtoCn8xuxfl0rwwas4U7mGM\nFqr0HwBFbE+RTsZeDTxozPipDI7n2NpvCraHQtDUc3XGr58iP/J3fJoy28gKQs4W\nB1dS+ZhsASpjUWMKg6BzYdcLa8tHR1bt9UBz5H25yZAzle3xz9KKI/rEEziWItoL\nSs/aeGFU96Xdo0BkvAjavoE47lR626AX+Q5BLgVRgh+pzBz8fviruKxPKTpOwoRn\neFkUq/4jqyHLxCl7ez4fpLT+UttjvYgHaMo+6WapaZl6HAGDw7ZEzb1rZKjPzanK\nMFIvAgMBAAGjgYMwgYAwHQYDVR0OBBYEFLjlGNlJb7o321Rrk2Bv+iSwYznKMB8G\nA1UdIwQYMBaAFLjlGNlJb7o321Rrk2Bv+iSwYznKMA8GA1UdEwEB/wQFMAMBAf8w\nLQYDVR0RBCYwJIILZXhhbXBsZS5jb22CD3d3dy5leGFtcGxlLm5ldIcECgAAATAN\nBgkqhkiG9w0BAQsFAAOCAgEAAvcpgKcHzu3HhA1Ew2IS7IXRUlYjV5XXs32Lxhns\nbQilYcYzG4Lc6ODBcbfMowO1T1gIHCl/AWNNq3+IE08fQlDkzYKqZupgnkAutU6B\nJpJc7sozjqvrxzbPLnI1sydY8YsVTFFhwELDfMbjqyQUE8Wtozq011A5VUuRI4ih\ngy7trKLFJg0IVmnRupeyzP1HgMxoEgqkWn/Y8AIL5VnFhVa6kwYdyux29/oUs+sy\npsJZs39YIDsntkId+T4/YPehW0GKe6gxPhUWUSCAOIeopEnIaL/us3b2tRzPjL1L\nu1sHG0O6xFakrjFgYeVH6hjBWjgNbQl01Gq3/jgM24XLK+fb2fNonDWA9eT6fc7L\nLKe5MoN0w5JSAMA9EBYdI7Y/rLq6CZ00Vzn/FXCoab/USWyhR2ESrF8L5u4mLNCZ\nEd1h8P5WwzPXuKWEf//QN75vqufS1z6jUBE9ratnpwgK0/FZkkotkvmad17/b8se\n/fIZ7vhVPCPPmE9NIy+hL976H2r5iqsBM9Tdcrfy9P4aQVtg0OoyKIPx5qbkgP34\n+og4VhwdbL2EAmNiVpyHWGpYTZn9UYCTgZ54/a+IvmH/5o8ieyKI19a6giasFUQQ\n1U53G7Vc5CF2bJJOBzXNn4qHGSut8hxdyiqIo1Ug//LYJRjxWuZaPxYkGnVpInYS\nYL0=\n-----END CERTIFICATE-----\n"
	private_key = "-----BEGIN PRIVATE KEY-----\n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\nPf6OYBciIGg0MzhUNspVIGRdT47WJQ==\n-----END PRIVATE KEY-----\n"
	labels = {
		"version" = "v1"
	}
}	
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_ssl_certificate.test", "id"),
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "type", "server"),
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "status", "1"),
					// Snis plan modifier checks
					// Verify snis count
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "snis.#", "2"),
					// Verify each value of the snis list
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "snis.0", "example.com"),
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "snis.1", "www.example.net"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_ssl_certificate.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore private_key during import as APISIX uses base64 format for it
				ImportStateVerifyIgnore: []string{"private_key"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_ssl_certificate" "test" {
	certificate = "-----BEGIN CERTIFICATE-----\nMIIFPjCCAyagAwIBAgIUe/c2A8XOSdHw60+OrCYfh3dzmm4wDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wHhcNMjMwNjE0MTAxOTMzWhcNMzMw\nNjExMTAxOTMzWjAWMRQwEgYDVQQDDAtleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcN\nAQEBBQADggIPADCCAgoCggIBALAxaEUndKUJ8whS48/pLT5Md/2LuQaov/0lpRcS\naPuteQq0lrbymZA9K0aAUt4Slx1zkKmfW4MKMHtsT5aq3uE7h2+I7v5mnT0EW6+1\nMXwmcLuqQrCVg55ogYHNcoLtsiSpd+vQU+F+MoV0iJpcmf0HcFKN1U9UNldkggl2\nQxIsvaPWxzEZMerHmgNVkAG7YlEXO20bgJfsbQeoBb0n8bNbr09nMcqlcV4hD0H+\nRKf3Xo4o9JdeIJGXbFf8JGolqkFSNUhlwPTkf4Ksndq0nALE0rqAVH848mTdZf58\n+JupCBCvwrajRGI4XKBYziAUe+0GVWIDSalzw6z7OSodOtLQvAIgMauGfJugLHRy\nF677CZIykVv4veAZat7RhI9XeHK83PMxYk+dQvAMOVtoCn8xuxfl0rwwas4U7mGM\nFqr0HwBFbE+RTsZeDTxozPipDI7n2NpvCraHQtDUc3XGr58iP/J3fJoy28gKQs4W\nB1dS+ZhsASpjUWMKg6BzYdcLa8tHR1bt9UBz5H25yZAzle3xz9KKI/rEEziWItoL\nSs/aeGFU96Xdo0BkvAjavoE47lR626AX+Q5BLgVRgh+pzBz8fviruKxPKTpOwoRn\neFkUq/4jqyHLxCl7ez4fpLT+UttjvYgHaMo+6WapaZl6HAGDw7ZEzb1rZKjPzanK\nMFIvAgMBAAGjgYMwgYAwHQYDVR0OBBYEFLjlGNlJb7o321Rrk2Bv+iSwYznKMB8G\nA1UdIwQYMBaAFLjlGNlJb7o321Rrk2Bv+iSwYznKMA8GA1UdEwEB/wQFMAMBAf8w\nLQYDVR0RBCYwJIILZXhhbXBsZS5jb22CD3d3dy5leGFtcGxlLm5ldIcECgAAATAN\nBgkqhkiG9w0BAQsFAAOCAgEAAvcpgKcHzu3HhA1Ew2IS7IXRUlYjV5XXs32Lxhns\nbQilYcYzG4Lc6ODBcbfMowO1T1gIHCl/AWNNq3+IE08fQlDkzYKqZupgnkAutU6B\nJpJc7sozjqvrxzbPLnI1sydY8YsVTFFhwELDfMbjqyQUE8Wtozq011A5VUuRI4ih\ngy7trKLFJg0IVmnRupeyzP1HgMxoEgqkWn/Y8AIL5VnFhVa6kwYdyux29/oUs+sy\npsJZs39YIDsntkId+T4/YPehW0GKe6gxPhUWUSCAOIeopEnIaL/us3b2tRzPjL1L\nu1sHG0O6xFakrjFgYeVH6hjBWjgNbQl01Gq3/jgM24XLK+fb2fNonDWA9eT6fc7L\nLKe5MoN0w5JSAMA9EBYdI7Y/rLq6CZ00Vzn/FXCoab/USWyhR2ESrF8L5u4mLNCZ\nEd1h8P5WwzPXuKWEf//QN75vqufS1z6jUBE9ratnpwgK0/FZkkotkvmad17/b8se\n/fIZ7vhVPCPPmE9NIy+hL976H2r5iqsBM9Tdcrfy9P4aQVtg0OoyKIPx5qbkgP34\n+og4VhwdbL2EAmNiVpyHWGpYTZn9UYCTgZ54/a+IvmH/5o8ieyKI19a6giasFUQQ\n1U53G7Vc5CF2bJJOBzXNn4qHGSut8hxdyiqIo1Ug//LYJRjxWuZaPxYkGnVpInYS\nYL0=\n-----END CERTIFICATE-----\n"
	private_key = "-----BEGIN PRIVATE KEY-----\n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\nPf6OYBciIGg0MzhUNspVIGRdT47WJQ==\n-----END PRIVATE KEY-----\n"
	labels = {
		"version" = "v2"
	}
	type = "client"
	snis = ["example.com"]
	status = 0
}	
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_ssl_certificate.test", "id"),
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "type", "client"),
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "status", "0"),
					// Snis plan modifier checks
					// Verify snis count
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "snis.#", "1"),
					// Verify each value of the snis list
					resource.TestCheckResourceAttr("apisix_ssl_certificate.test", "snis.0", "example.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
